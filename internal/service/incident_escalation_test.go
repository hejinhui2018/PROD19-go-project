package service

import (
	"testing"
	"time"

	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/domain"
)

func TestCriticalIncidentEscalationReachesControlCenter(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 8, 25, 3, 40, 0, 0, time.UTC)
	s, err := New(dir, clock.Fixed{T: now})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	blockade := domain.Blockade{
		ID:      "RG-684",
		Name:    "overnight obstruction response",
		Segment: domain.TrackSegment{ID: "seg-east", Name: "East Main"},
		Crew:    domain.Crew{ID: "crew-east", Name: "east crew", Qualified: true},
		Locks: domain.SafetyLockPack{
			ID:        "lock-684",
			Applied:   true,
			AppliedAt: now,
		},
		Approval: domain.DispatcherApproval{
			ID:         "appr-684",
			Dispatcher: "dispatcher-east",
			Approved:   true,
			At:         now,
		},
		Equipment: domain.EquipmentHandoff{
			ID:        "handoff-684",
			Custodian: "crew-east",
			Complete:  true,
			At:        now,
		},
		Start: now,
		End:   now.Add(90 * time.Minute),
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

	workflow := NewWorkflow(s, "lead-east", now)
	if _, err := workflow.Incident(blockade.ID, "foreign object near protected track"); err != nil {
		t.Fatalf("record incident: %v", err)
	}

	restarted, err := New(dir, clock.Fixed{T: now.Add(3 * time.Minute)})
	if err != nil {
		t.Fatalf("restart service: %v", err)
	}
	if _, err := restarted.RecoveryReport(); err != nil {
		t.Fatalf("recovery report: %v", err)
	}
	recovered, err := restarted.Get(blockade.ID)
	if err != nil {
		t.Fatalf("get recovered blockade: %v", err)
	}
	if recovered.State != domain.StateIncident {
		t.Fatalf("state = %s, want %s", recovered.State, domain.StateIncident)
	}

	actions, err := restarted.EscalationList(blockade.ID)
	if err != nil {
		t.Fatalf("escalation list: %v", err)
	}
	teams := map[string]bool{}
	for _, action := range actions {
		teams[action.Team] = true
	}
	for _, team := range []string{"dispatcher", "signal-desk", "control-center"} {
		if !teams[team] {
			t.Fatalf("missing escalation team %q in %#v", team, actions)
		}
	}
}
