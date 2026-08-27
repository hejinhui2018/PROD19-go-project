package recovery

import (
	"github.com/wogo-prod19/railguard/internal/domain"
	"time"
)

type RestartPlan struct {
	Snapshot domain.Blockade
	Events   []domain.Event
	Started  time.Time
}

func (p RestartPlan) Apply() domain.Blockade {
	b := p.Snapshot
	for _, e := range p.Events {
		if e.BlockadeID == b.ID && e.Payload.Version >= b.Version {
			b = e.Payload
		}
	}
	return b
}
func (p RestartPlan) Complete(now time.Time) RestartReport {
	r, _ := VerifyReplay(p.Events)
	return RestartReport{StartedAt: p.Started, CompletedAt: now, Replay: r}
}
