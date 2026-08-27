package safety

import "time"

type EscalationStep struct {
	Team      string
	Deadline  time.Time
	Completed bool
}
type EscalationPlan struct {
	IncidentID string
	Steps      []EscalationStep
}

func (p *EscalationPlan) Next(now time.Time) *EscalationStep {
	for i := range p.Steps {
		if !p.Steps[i].Completed && now.After(p.Steps[i].Deadline) {
			return &p.Steps[i]
		}
	}
	return nil
}
func (p *EscalationPlan) Complete(team string) bool {
	for i := range p.Steps {
		if p.Steps[i].Team == team && !p.Steps[i].Completed {
			p.Steps[i].Completed = true
			return true
		}
	}
	return false
}
