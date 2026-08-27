package service

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/domain"
	"github.com/wogo-prod19/railguard/internal/planning"
	"github.com/wogo-prod19/railguard/internal/safety"
	"time"
)

type Workflow struct {
	Service *Service
	Actor   string
	Started time.Time
}

func NewWorkflow(s *Service, actor string, now time.Time) Workflow { return Workflow{s, actor, now} }
func (w Workflow) Draft(id string) (planning.Draft, error)         { return w.Service.Plan(id, w.Actor) }
func (w Workflow) Prepare(id string) (planning.Draft, error) {
	d, e := w.Draft(id)
	if e != nil {
		return d, e
	}
	if issues := planning.CheckAdjacent(d.Segments, d.Dependencies); len(issues) > 0 {
		return d, fmt.Errorf("dependencies unresolved")
	}
	if !d.Ready() {
		return d, fmt.Errorf("readiness incomplete")
	}
	return d, nil
}
func (w Workflow) Protect(id string, lock safety.LockCheck) (domain.Blockade, error) {
	b, e := w.Service.Get(id)
	if e != nil {
		return b, e
	}
	if e = safety.VerifyLockPack(b.Locks, lock); e != nil {
		return b, e
	}
	if e = w.Service.VerifySafety(id); e != nil {
		return b, e
	}
	if e = w.Service.Transition(id, domain.StateProtected); e != nil {
		return b, e
	}
	return w.Service.Get(id)
}
func (w Workflow) Activate(id string) (domain.Blockade, error) {
	b, e := w.Service.Get(id)
	if e != nil {
		return b, e
	}
	if !ReadyForProtection(b) {
		return b, fmt.Errorf("protection prerequisites incomplete")
	}
	if e = w.Service.Transition(id, domain.StateActive); e != nil {
		return b, e
	}
	return w.Service.Get(id)
}
func (w Workflow) RequestRelease(id string) (domain.Blockade, error) {
	if e := w.Service.Transition(id, domain.StateReleasePending); e != nil {
		return domain.Blockade{}, e
	}
	return w.Service.Get(id)
}
func (w Workflow) Release(id string) (domain.Blockade, error) {
	ok, _, e := w.Service.ReleaseReady(id)
	if e != nil {
		return domain.Blockade{}, e
	}
	if !ok {
		return domain.Blockade{}, fmt.Errorf("release compliance checks failed")
	}
	if e = w.Service.Transition(id, domain.StateReleased); e != nil {
		return domain.Blockade{}, e
	}
	return w.Service.Get(id)
}
func (w Workflow) Incident(id string, summary string) (domain.Blockade, error) {
	return w.RegisterIncident(id, "critical", summary)
}
func (w Workflow) RegisterIncident(id, severity, summary string) (domain.Blockade, error) {
	if summary == "" {
		return domain.Blockade{}, fmt.Errorf("incident summary required")
	}
	return w.Service.RegisterIncident(id, severity, w.Actor, summary)
}
func (w Workflow) Observe(id, msg string) (domain.Blockade, error) {
	if msg == "" {
		return domain.Blockade{}, fmt.Errorf("observation required")
	}
	e := w.Service.AddObservation(id, domain.FieldObservation{ID: fmt.Sprintf("obs-%d", w.Started.UnixNano()), Reporter: w.Actor, Message: msg, At: time.Now()})
	if e != nil {
		return domain.Blockade{}, e
	}
	return w.Service.Get(id)
}
func (w Workflow) Diagnostics() string { return w.Service.Diagnostics().Summary() }
