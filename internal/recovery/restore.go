package recovery

import "github.com/wogo-prod19/railguard/internal/domain"

type RestoreResult struct {
	Blockade domain.Blockade
	Applied  int
	Skipped  int
	Warnings []string
}

func Restore(base domain.Blockade, events []domain.Event) RestoreResult {
	r := RestoreResult{Blockade: base}
	for _, e := range events {
		if e.BlockadeID != base.ID {
			r.Skipped++
			continue
		}
		if e.Payload.Version > r.Blockade.Version {
			r.Blockade = e.Payload
			r.Applied++
		} else {
			r.Warnings = append(r.Warnings, "stale event")
		}
	}
	return r
}
func (r RestoreResult) Successful() bool { return r.Applied > 0 && len(r.Warnings) == 0 }
