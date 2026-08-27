package recovery

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/domain"
)

type ReplayReport struct {
	Events     int
	LastNumber int64
	Blockades  int
	Errors     []string
}

func VerifyReplay(events []domain.Event) (ReplayReport, error) {
	r := ReplayReport{}
	seen := map[string]bool{}
	var last int64
	for _, e := range events {
		if e.Type == domain.EventObservation {
			continue
		}
		r.Events++
		if e.Number <= last {
			r.Errors = append(r.Errors, fmt.Sprintf("event order at %d", e.Number))
		}
		last = e.Number
		if e.BlockadeID == "" {
			r.Errors = append(r.Errors, "event missing blockade")
		}
		seen[e.BlockadeID] = true
	}
	r.LastNumber = last
	r.Blockades = len(seen)
	if len(r.Errors) > 0 {
		return r, fmt.Errorf("replay verification failed")
	}
	return r, nil
}
