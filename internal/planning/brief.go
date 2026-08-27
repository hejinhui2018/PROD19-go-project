package planning

import "strings"

type Brief struct {
	Subject      string
	Objectives   []string
	Constraints  []string
	Acknowledged map[string]bool
}

func (b Brief) Valid() bool { return strings.TrimSpace(b.Subject) != "" && len(b.Objectives) > 0 }
func (b *Brief) Acknowledge(actor string) {
	if b.Acknowledged == nil {
		b.Acknowledged = map[string]bool{}
	}
	b.Acknowledged[actor] = true
}
func (b Brief) FullyAcknowledged(actors []string) bool {
	for _, a := range actors {
		if !b.Acknowledged[a] {
			return false
		}
	}
	return true
}
