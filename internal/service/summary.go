package service

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/domain"
	"strings"
)

func (s *Service) OperationalSummary() string {
	counts := s.StateCounts()
	parts := []string{}
	for _, st := range []domain.State{domain.StateDraft, domain.StatePlanned, domain.StateCrewReady, domain.StateProtected, domain.StateActive, domain.StateReleasePending, domain.StateReleased, domain.StateIncident} {
		if n := counts[st]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", st, n))
		}
	}
	if len(parts) == 0 {
		return "no blockades"
	}
	return strings.Join(parts, ",")
}
func (s *Service) HasState(st domain.State) bool { return s.StateCounts()[st] > 0 }
func (s *Service) ActiveCount() int              { return s.StateCounts()[domain.StateActive] }
func (s *Service) IncidentCount() int            { return s.StateCounts()[domain.StateIncident] }
func (s *Service) ReleasePendingCount() int      { return s.StateCounts()[domain.StateReleasePending] }
