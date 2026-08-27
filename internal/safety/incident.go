package safety

import (
	"fmt"
	"time"
)

type Incident struct {
	ID, BlockadeID, Severity, Reporter, Summary string
	OpenedAt                                    time.Time
	Escalated                                   bool
	Escalation                                  []string
}

func Escalate(i Incident, now time.Time) (Incident, error) {
	if i.ID == "" || i.BlockadeID == "" || i.Summary == "" {
		return i, fmt.Errorf("incident incomplete")
	}
	i.OpenedAt = now
	i.Escalated = true
	if i.Severity == "critical" {
		i.Escalation = []string{"dispatcher", "signal-desk", "control-center"}
	} else {
		i.Escalation = []string{"dispatcher", "signal-desk"}
	}
	return i, nil
}
