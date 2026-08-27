package service

import (
	"github.com/wogo-prod19/railguard/internal/domain"
	"github.com/wogo-prod19/railguard/internal/recovery"
)

func (s *Service) Diagnostics() recovery.DiagnosticSet {
	d := recovery.DiagnosticSet{}
	ev, e := s.events.Load()
	if e != nil {
		d.Add(recovery.Diagnostic{Component: "event-store", Status: "error", Detail: e.Error(), Critical: true})
	} else {
		_, e = recovery.VerifyReplay(ev)
		if e != nil {
			d.Add(recovery.Diagnostic{Component: "replay", Status: "error", Detail: e.Error(), Critical: true})
		} else {
			d.Add(recovery.Diagnostic{Component: "replay", Status: "ok", Detail: "events verified", Critical: true})
		}
	}
	if len(s.blocks) == 0 {
		d.Add(recovery.Diagnostic{Component: "registry", Status: "warning", Detail: "no active blockade", Critical: false})
	} else {
		d.Add(recovery.Diagnostic{Component: "registry", Status: "ok", Detail: "blockades loaded", Critical: true})
	}
	_ = domain.StateDraft
	return d
}
