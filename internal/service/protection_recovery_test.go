package service_test

import (
	"testing"
	"time"

	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/domain"
	"github.com/wogo-prod19/railguard/internal/safety"
	"github.com/wogo-prod19/railguard/internal/service"
)

func TestFieldProtectionRejectionRecovery(t *testing.T) {
	now := time.Date(2026, 8, 25, 3, 40, 0, 0, time.UTC)
	dir := t.TempDir()
	svc, err := service.New(dir, clock.Fixed{T: now})
	if err != nil {
		t.Fatalf("open service: %v", err)
	}

	blockade := domain.Blockade{
		ID:      "RG-411",
		Name:    "overnight turnout protection",
		Segment: domain.TrackSegment{ID: "SEG-411", Name: "north throat"},
		Crew:    domain.Crew{ID: "crew-night", Name: "night possession", Qualified: true},
		Locks: domain.SafetyLockPack{
			ID:        "lock-RG-411",
			Applied:   true,
			AppliedAt: now,
		},
		Approval: domain.DispatcherApproval{
			ID:         "approval-RG-411",
			Dispatcher: "dispatcher-17",
			Approved:   true,
			At:         now,
		},
		Equipment: domain.EquipmentHandoff{
			ID:        "kit-RG-411",
			Custodian: "crew-night",
			Complete:  true,
			At:        now,
		},
		Start: now.Add(time.Hour),
		End:   now.Add(3 * time.Hour),
	}
	if err := svc.Create(blockade); err != nil {
		t.Fatalf("create blockade: %v", err)
	}
	for _, state := range []domain.State{domain.StatePlanned, domain.StateCrewReady} {
		if err := svc.Transition(blockade.ID, state); err != nil {
			t.Fatalf("move to %s: %v", state, err)
		}
	}
	before, err := svc.RecoveryReport()
	if err != nil {
		t.Fatalf("pre-check recovery report: %v", err)
	}

	workflow := service.NewWorkflow(svc, "dispatcher-17", now)
	_, err = workflow.Protect(blockade.ID, safety.LockCheck{
		LockID:     "lock-RG-411",
		Inspector:  "field-lead",
		VerifiedAt: now,
		Evidence:   nil,
	})
	if err == nil {
		t.Fatalf("field protection unexpectedly succeeded")
	}

	current, err := svc.Get(blockade.ID)
	if err != nil {
		t.Fatalf("read current blockade: %v", err)
	}
	if current.State != domain.StateCrewReady {
		t.Fatalf("current state after rejected protection = %s", current.State)
	}

	reopened, err := service.New(dir, clock.Fixed{T: now.Add(time.Minute)})
	if err != nil {
		t.Fatalf("reopen service: %v", err)
	}
	recovered, err := reopened.Get(blockade.ID)
	if err != nil {
		t.Fatalf("read recovered blockade: %v", err)
	}
	if recovered.State != domain.StateCrewReady {
		t.Fatalf("recovered state after rejected protection = %s", recovered.State)
	}
	after, err := reopened.RecoveryReport()
	if err != nil {
		t.Fatalf("post-check recovery report: %v", err)
	}
	if after.Replay.Events != before.Replay.Events {
		t.Fatalf("recovery event count changed after rejected protection: before=%d after=%d", before.Replay.Events, after.Replay.Events)
	}
}
