package safety

import "time"

type Briefing struct {
	ID, Leader   string
	Topics       []string
	Attendees    []string
	At           time.Time
	Acknowledged map[string]bool
}

func (b Briefing) Ready() bool {
	if b.ID == "" || b.Leader == "" || len(b.Topics) == 0 || len(b.Attendees) == 0 {
		return false
	}
	for _, a := range b.Attendees {
		if !b.Acknowledged[a] {
			return false
		}
	}
	return true
}
func (b *Briefing) Acknowledge(name string) {
	if b.Acknowledged == nil {
		b.Acknowledged = map[string]bool{}
	}
	b.Acknowledged[name] = true
}
