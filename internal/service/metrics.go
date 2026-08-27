package service

import (
	"github.com/wogo-prod19/railguard/internal/domain"
	"sync"
	"time"
)

type Metrics struct {
	mu                                            sync.Mutex
	Created, Transitions, Observations, Incidents int
	LastTransition                                time.Time
}

func (m *Metrics) RecordCreate() { m.mu.Lock(); defer m.mu.Unlock(); m.Created++ }
func (m *Metrics) RecordTransition(at time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Transitions++
	m.LastTransition = at
}
func (m *Metrics) RecordObservation() { m.mu.Lock(); defer m.mu.Unlock(); m.Observations++ }
func (m *Metrics) RecordIncident()    { m.mu.Lock(); defer m.mu.Unlock(); m.Incidents++ }

type MetricsSnapshot struct {
	Created, Transitions, Observations, Incidents int
	LastTransition                                time.Time
}

func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	return MetricsSnapshot{m.Created, m.Transitions, m.Observations, m.Incidents, m.LastTransition}
}

type Health struct {
	StateCounts   map[domain.State]int
	NoticePending int
	AuditRecords  int
	Ready         bool
}

func (s *Service) Health() Health {
	counts := s.StateCounts()
	ns, _ := s.Outbox()
	ok, _, _ := s.ReleaseReadyForAny()
	return Health{counts, len(ns), len(s.AuditRecords()), ok}
}
func (s *Service) ReleaseReadyForAny() (bool, domain.Blockade, error) {
	for _, b := range s.List() {
		ok, _, e := s.ReleaseReady(b.ID)
		if e == nil && ok {
			return true, b, nil
		}
	}
	return false, domain.Blockade{}, nil
}
