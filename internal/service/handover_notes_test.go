package service

import (
	"testing"
	"time"

	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/domain"
)

func TestRestartKeepsHandoverNotes(t *testing.T) {
	now := time.Date(2026, 8, 25, 22, 30, 0, 0, time.UTC)
	dir := t.TempDir()
	svc, err := New(dir, clock.Fixed{T: now})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	blockade := domain.Blockade{
		ID:      "RG-815",
		Name:    "north lock review",
		Segment: domain.TrackSegment{ID: "N-17", Name: "north approach"},
		Crew:    domain.Crew{ID: "crew-night", Name: "night maintenance", Qualified: true},
		Start:   now,
		End:     now.Add(90 * time.Minute),
	}
	if err := svc.Create(blockade); err != nil {
		t.Fatalf("create blockade: %v", err)
	}
	for _, state := range []domain.State{domain.StatePlanned, domain.StateCrewReady, domain.StateProtected} {
		if err := svc.Transition(blockade.ID, state); err != nil {
			t.Fatalf("transition to %s: %v", state, err)
		}
	}
	note := domain.FieldObservation{
		ID:       "obs-north-lock",
		Reporter: "field-lead",
		Message:  "north end guard rechecked the lock kit",
		At:       now.Add(35 * time.Minute),
	}
	if err := svc.AddObservation(blockade.ID, note); err != nil {
		t.Fatalf("add handover note: %v", err)
	}
	live, err := svc.Get(blockade.ID)
	if err != nil {
		t.Fatalf("get live blockade: %v", err)
	}
	if len(live.Observations) != 1 {
		t.Fatalf("live handover note count = %d", len(live.Observations))
	}
	recovered, err := New(dir, clock.Fixed{T: now.Add(time.Hour)})
	if err != nil {
		t.Fatalf("restart service: %v", err)
	}
	afterRestart, err := recovered.Get(blockade.ID)
	if err != nil {
		t.Fatalf("get recovered blockade: %v", err)
	}
	if len(afterRestart.Observations) != 1 {
		t.Fatalf("recovered handover note count = %d", len(afterRestart.Observations))
	}
	if afterRestart.Observations[0].Message != note.Message {
		t.Fatalf("recovered handover note message = %q", afterRestart.Observations[0].Message)
	}
	report, err := recovered.RecoveryReport()
	if err != nil {
		t.Fatalf("recovery report: %v", err)
	}
	if report.Replay.Events != 5 {
		t.Fatalf("recovery event count = %d", report.Replay.Events)
	}
}
