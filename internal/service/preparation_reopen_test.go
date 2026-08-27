package service_test

import (
	"testing"
	"time"

	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/domain"
	"github.com/wogo-prod19/railguard/internal/service"
)

func TestPreparationDraftIncludesStoredCoordination(t *testing.T) {
	now := time.Date(2026, 8, 25, 8, 30, 0, 0, time.UTC)
	dir := t.TempDir()
	svc, err := service.New(dir, clock.Fixed{T: now})
	if err != nil {
		t.Fatalf("open service: %v", err)
	}

	blockade := domain.Blockade{
		ID:      "RG-926",
		Name:    "morning turnout possession",
		Segment: domain.TrackSegment{ID: "ML-926", Name: "main line throat"},
		Crew:    domain.Crew{ID: "crew-main", Name: "main line team", Qualified: true},
		Locks: domain.SafetyLockPack{
			ID:        "lock-RG-926",
			Applied:   true,
			AppliedAt: now,
		},
		Approval: domain.DispatcherApproval{
			ID:         "approval-RG-926",
			Dispatcher: "dispatcher-11",
			Approved:   true,
			At:         now,
		},
		Equipment: domain.EquipmentHandoff{
			ID:        "kit-RG-926",
			Custodian: "crew-main",
			Complete:  true,
			At:        now,
		},
		Start: now.Add(30 * time.Minute),
		End:   now.Add(2 * time.Hour),
		Adjacent: []domain.SegmentLink{{
			From:     "ML-926",
			To:       "AL-212",
			Distance: 1,
		}},
	}
	if err := svc.Create(blockade); err != nil {
		t.Fatalf("create blockade: %v", err)
	}

	reopened, err := service.New(dir, clock.Fixed{T: now.Add(time.Minute)})
	if err != nil {
		t.Fatalf("reopen service: %v", err)
	}
	current, err := reopened.Get(blockade.ID)
	if err != nil {
		t.Fatalf("read reopened blockade: %v", err)
	}
	if len(current.Adjacent) != 1 {
		t.Fatalf("reopened coordination count = %d", len(current.Adjacent))
	}

	workflow := service.NewWorkflow(reopened, "dispatcher-11", now.Add(time.Minute))
	draft, err := workflow.Prepare(blockade.ID)
	if err != nil {
		t.Fatalf("prepare after reopen: %v", err)
	}
	if !containsSegment(draft.Segments, "AL-212") {
		t.Fatalf("preparation segments = %#v, want adjacent segment", draft.Segments)
	}
	if len(draft.Dependencies) != 1 || draft.Dependencies[0].To != "AL-212" {
		t.Fatalf("preparation dependencies = %#v", draft.Dependencies)
	}
}

func containsSegment(segments []domain.TrackSegment, id string) bool {
	for _, segment := range segments {
		if segment.ID == id {
			return true
		}
	}
	return false
}
