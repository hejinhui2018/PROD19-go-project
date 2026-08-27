package safety

import "time"

type EscalationStep struct {
	Team      string
	Deadline  time.Time
	Completed bool
}
type EscalationPlan struct {
	IncidentID string
	Immediate  []EscalationStep
	FollowUp   []EscalationStep
}

func (p *EscalationPlan) Next(now time.Time) *EscalationStep {
	for i := range p.Immediate {
		if !p.Immediate[i].Completed && now.After(p.Immediate[i].Deadline) {
			return &p.Immediate[i]
		}
	}
	for i := range p.FollowUp {
		if !p.FollowUp[i].Completed && now.After(p.FollowUp[i].Deadline) {
			return &p.FollowUp[i]
		}
	}
	return nil
}
func (p *EscalationPlan) Complete(team string) bool {
	for _, steps := range []*[]EscalationStep{&p.Immediate, &p.FollowUp} {
		for i := range *steps {
			if (*steps)[i].Team == team && !(*steps)[i].Completed {
				(*steps)[i].Completed = true
				return true
			}
		}
	}
	return false
}

func BuildEscalationPlan(i Incident) EscalationPlan {
	p := EscalationPlan{IncidentID: i.ID}
	for _, team := range i.Escalation {
		step := EscalationStep{Team: team, Deadline: i.OpenedAt.Add(5 * time.Minute)}
		if team == "control-center" {
			step.Deadline = i.OpenedAt.Add(10 * time.Minute)
			p.FollowUp = append(p.FollowUp, step)
			continue
		}
		p.Immediate = append(p.Immediate, step)
	}
	return p
}
