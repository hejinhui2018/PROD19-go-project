package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/domain"
)

func TestArchivedLogKeepsActiveBlockadeInHandoff(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 8, 25, 2, 15, 0, 0, time.UTC)
	s, err := New(dir, clock.Fixed{T: now})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	blockade := domain.Blockade{
		ID:      "RG-538",
		Name:    "night rail replacement",
		Segment: domain.TrackSegment{ID: "seg-night", Name: "Night Up Main"},
		Crew:    domain.Crew{ID: "crew-night", Name: "night crew", Qualified: true},
		Locks: domain.SafetyLockPack{
			ID:        "lock-538",
			Applied:   true,
			AppliedAt: now,
		},
		Approval: domain.DispatcherApproval{
			ID:         "appr-538",
			Dispatcher: "dispatcher-night",
			Approved:   true,
			At:         now,
		},
		Equipment: domain.EquipmentHandoff{
			ID:        "handoff-538",
			Custodian: "crew-night",
			Complete:  true,
			At:        now,
		},
		Start: now,
		End:   now.Add(2 * time.Hour),
	}
	if err := s.Create(blockade); err != nil {
		t.Fatalf("create blockade: %v", err)
	}
	for _, state := range []domain.State{
		domain.StatePlanned,
		domain.StateCrewReady,
		domain.StateProtected,
		domain.StateActive,
	} {
		if err := s.Transition(blockade.ID, state); err != nil {
			t.Fatalf("transition to %s: %v", state, err)
		}
	}

	if err := os.Remove(filepath.Join(dir, "events.jsonl")); err != nil {
		t.Fatalf("archive event log: %v", err)
	}

	restarted, err := New(dir, clock.Fixed{T: now.Add(5 * time.Minute)})
	if err != nil {
		t.Fatalf("restart service: %v", err)
	}
	if _, err := restarted.RecoveryReport(); err != nil {
		t.Fatalf("recovery report: %v", err)
	}

	items := restarted.List()
	if len(items) != 1 {
		t.Fatalf("handoff list count = %d, want 1", len(items))
	}
	if items[0].ID != blockade.ID || items[0].State != domain.StateActive {
		t.Fatalf("handoff item = %s/%s, want %s/%s", items[0].ID, items[0].State, blockade.ID, domain.StateActive)
	}
	if summary := restarted.OperationalSummary(); !strings.Contains(summary, "active=1") {
		t.Fatalf("summary = %q, want active blockade count", summary)
	}
}
