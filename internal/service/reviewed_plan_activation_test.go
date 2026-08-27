package service

import (
	"strings"
	"testing"
	"time"

	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/domain"
)

func TestReviewedPlanActivationAfterRestart(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 8, 25, 14, 30, 0, 0, time.UTC)
	s, err := New(dir, clock.Fixed{T: now})
	if err != nil {
		t.Fatal(err)
	}
	blockade := domain.Blockade{
		ID:      "RG-759",
		Name:    "afternoon possession",
		Segment: domain.TrackSegment{ID: "seg-759", Name: "south turnout"},
		Crew:    domain.Crew{ID: "crew-759", Name: "field team", Qualified: true},
		Locks: domain.SafetyLockPack{
			ID:        "lock-759",
			Applied:   true,
			AppliedAt: now,
		},
		Approval: domain.DispatcherApproval{
			ID:         "approval-759",
			Dispatcher: "dispatcher-a",
			Approved:   true,
			At:         now,
		},
		Equipment: domain.EquipmentHandoff{
			ID:        "handoff-759",
			Custodian: "crew-759",
			Complete:  true,
			At:        now,
		},
		Start: now.Add(30 * time.Minute),
		End:   now.Add(2 * time.Hour),
	}
	if err := s.Create(blockade); err != nil {
		t.Fatal(err)
	}
	for _, state := range []domain.State{domain.StatePlanned, domain.StateCrewReady, domain.StateProtected} {
		if err := s.Transition("RG-759", state); err != nil {
			t.Fatalf("transition to %s: %v", state, err)
		}
	}
	w := NewWorkflow(s, "dispatcher-a", now)
	if _, err := w.ReviewPlan("RG-759", 2); err != nil {
		t.Fatal(err)
	}
	restarted, err := New(dir, clock.Fixed{T: now.Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := restarted.Get("RG-759")
	if err != nil {
		t.Fatal(err)
	}
	if recovered.PlanRevision != 2 {
		t.Fatalf("expected recovered plan revision 2, got %d", recovered.PlanRevision)
	}
	workflow := NewWorkflow(restarted, "field-lead", now.Add(time.Minute))
	if _, err := workflow.Prepare("RG-759"); err != nil {
		t.Fatalf("prepare after restart: %v", err)
	}
	activated, err := workflow.Activate("RG-759")
	if err != nil {
		t.Fatalf("activate after reviewed plan restart: %v", err)
	}
	if activated.State != domain.StateActive {
		t.Fatalf("expected active state, got %s", activated.State)
	}
	if summary := restarted.OperationalSummary(); !strings.Contains(summary, "active=1") {
		t.Fatalf("expected active summary, got %s", summary)
	}
}
