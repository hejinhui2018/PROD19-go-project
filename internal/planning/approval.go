package planning

import (
	"fmt"
	"time"
)

type Approval struct {
	ID, Actor, Decision, Comment string
	At                           time.Time
}
type ApprovalTrail struct{ Entries []Approval }

func (a *ApprovalTrail) Record(x Approval) error {
	if x.ID == "" || x.Actor == "" || x.Decision == "" {
		return fmt.Errorf("approval incomplete")
	}
	a.Entries = append(a.Entries, x)
	return nil
}
func (a ApprovalTrail) Latest() Approval {
	if len(a.Entries) == 0 {
		return Approval{}
	}
	return a.Entries[len(a.Entries)-1]
}
func (a ApprovalTrail) Accepted() bool { return a.Latest().Decision == "approve" }
func (a ApprovalTrail) Actors() []string {
	out := []string{}
	seen := map[string]bool{}
	for _, x := range a.Entries {
		if !seen[x.Actor] {
			seen[x.Actor] = true
			out = append(out, x.Actor)
		}
	}
	return out
}
