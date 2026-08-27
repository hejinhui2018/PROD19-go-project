package safety

import (
	"fmt"
	"time"
)

type ComplianceWindow struct {
	Start, End    time.Time
	RequiredRoles []Role
}

func (c ComplianceWindow) Active(now time.Time) bool {
	return !now.Before(c.Start) && now.Before(c.End)
}
func (c ComplianceWindow) Validate() error {
	if !c.Start.Before(c.End) {
		return fmt.Errorf("invalid compliance window")
	}
	if len(c.RequiredRoles) == 0 {
		return fmt.Errorf("roles required")
	}
	return nil
}
func (c ComplianceWindow) Missing(r Roster) []Role {
	out := []Role{}
	for _, role := range c.RequiredRoles {
		if len(r.Members[role]) == 0 {
			out = append(out, role)
		}
	}
	return out
}

type AuthorityCheck struct {
	Name, Token string
	ValidUntil  time.Time
}

func (a AuthorityCheck) Valid(now time.Time) bool {
	return a.Name != "" && a.Token != "" && now.Before(a.ValidUntil)
}
