package audit

import "fmt"

type Check struct {
	Name   string
	Passed bool
	Detail string
}

func ReleaseChecks(records []Record, blockadeID string) []Check {
	seen := false
	for _, r := range records {
		if r.BlockadeID == blockadeID && r.To == "released" {
			seen = true
		}
	}
	return []Check{{"immutable-ledger", len(records) > 0, "audit records present"}, {"release-transition", seen, fmt.Sprintf("release recorded=%t", seen)}}
}
func Ready(checks []Check) bool {
	for _, c := range checks {
		if !c.Passed {
			return false
		}
	}
	return true
}
