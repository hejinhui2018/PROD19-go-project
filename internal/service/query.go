package service

import (
	"github.com/wogo-prod19/railguard/internal/audit"
	"github.com/wogo-prod19/railguard/internal/domain"
	"sort"
)

func (s *Service) List() []domain.Blockade {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.Blockade, 0, len(s.blocks))
	for _, b := range s.blocks {
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func (s *Service) AuditRecords() []audit.Record { return s.audit.List() }
func (s *Service) StateCounts() map[domain.State]int {
	out := map[domain.State]int{}
	for _, b := range s.List() {
		out[b.State]++
	}
	return out
}
