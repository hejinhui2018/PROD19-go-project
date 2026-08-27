package notice

import "time"

type DrainResult struct {
	Seen, Delivered, Retried, Failed int
	FinishedAt                       time.Time
	Errors                           []string
}

func (r DrainResult) Complete() bool     { return r.Failed == 0 && len(r.Errors) == 0 }
func NewDrain(now time.Time) DrainResult { return DrainResult{FinishedAt: now} }
