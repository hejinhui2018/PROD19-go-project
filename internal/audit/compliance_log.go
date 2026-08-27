package audit

import "time"

type ComplianceLog struct {
	Checks []Check
	Actor  string
	At     time.Time
}

func (l *ComplianceLog) Add(c Check) { l.Checks = append(l.Checks, c) }
func (l ComplianceLog) Passed() bool { return Ready(l.Checks) }
func (l ComplianceLog) FailedNames() []string {
	out := []string{}
	for _, c := range l.Checks {
		if !c.Passed {
			out = append(out, c.Name)
		}
	}
	return out
}

type ActorSummary struct {
	Actor   string
	Actions int
	Last    time.Time
}

func SummarizeActors(records []Record) []ActorSummary {
	m := map[string]*ActorSummary{}
	for _, r := range records {
		x := m[r.Actor]
		if x == nil {
			x = &ActorSummary{Actor: r.Actor}
			m[r.Actor] = x
		}
		x.Actions++
		if r.At.After(x.Last) {
			x.Last = r.At
		}
	}
	out := []ActorSummary{}
	for _, x := range m {
		out = append(out, *x)
	}
	return out
}
