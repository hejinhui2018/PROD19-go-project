package planning

import (
	"fmt"
	"time"
)

type RevisionResult struct {
	Draft    Draft
	Changed  bool
	Warnings []string
}

func Revise(previous Draft, next Draft, now time.Time) (RevisionResult, error) {
	if err := next.Validate(); err != nil {
		return RevisionResult{}, err
	}
	if previous.ID != "" && previous.ID != next.ID {
		return RevisionResult{}, fmt.Errorf("draft identity changed")
	}
	if next.Revision <= previous.Revision {
		return RevisionResult{}, fmt.Errorf("revision must advance")
	}
	next.UpdatedAt = now
	warnings := []string{}
	if len(next.Segments) != len(previous.Segments) {
		warnings = append(warnings, "segment topology changed")
	}
	return RevisionResult{next, true, warnings}, nil
}
func TransitionSafe(d Draft, active bool) error {
	if active && d.Revision < 2 {
		return fmt.Errorf("active blockade requires reviewed revision")
	}
	if len(CheckAdjacent(d.Segments, d.Dependencies)) > 0 {
		return fmt.Errorf("dependency checks failed")
	}
	return nil
}
