package service

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/audit"
	"github.com/wogo-prod19/railguard/internal/domain"
)

func (s *Service) ReleaseReady(id string) (bool, []audit.Check, error) {
	b, e := s.Get(id)
	if e != nil {
		return false, nil, e
	}
	checks := audit.ReleaseChecks(s.audit.List(), id)
	if b.State != domain.StateReleasePending {
		checks = append(checks, audit.Check{Name: "release-state", Passed: false, Detail: fmt.Sprintf("state=%s", b.State)})
	}
	return audit.Ready(checks), checks, nil
}
