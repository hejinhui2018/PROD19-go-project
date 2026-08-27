package service

import (
	"strings"
	"testing"
	"time"

	"github.com/wogo-prod19/railguard/internal/domain"
)

type advancingClock struct {
	next time.Time
	step time.Duration
}

func (c *advancingClock) Now() time.Time {
	c.next = c.next.Add(c.step)
	return c.next
}

func TestReleaseCompletionNotice(t *testing.T) {
	clk := &advancingClock{next: time.Date(2026, 8, 25, 21, 0, 0, 0, time.UTC), step: time.Second}
	dir := t.TempDir()
	svc, err := New(dir, clk)
	if err != nil {
		t.Fatal(err)
	}

	blockade := domain.Blockade{
		ID:      "RG-742",
		Name:    "northbound possession",
		Segment: domain.TrackSegment{ID: "SEG-N4", Name: "North 4"},
		Crew:    domain.Crew{ID: "crew-night", Name: "night crew", Qualified: true},
		Start:   clk.next,
		End:     clk.next.Add(2 * time.Hour),
	}
	if err := svc.Create(blockade); err != nil {
		t.Fatal(err)
	}
	for _, state := range []domain.State{
		domain.StatePlanned,
		domain.StateCrewReady,
		domain.StateProtected,
		domain.StateActive,
		domain.StateReleasePending,
		domain.StateReleased,
	} {
		if err := svc.Transition(blockade.ID, state); err != nil {
			t.Fatalf("transition to %s: %v", state, err)
		}
	}

	recovered, err := New(dir, clk)
	if err != nil {
		t.Fatal(err)
	}
	got, err := recovered.Get(blockade.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.State != domain.StateReleased {
		t.Fatalf("state page has %s", got.State)
	}

	notices, err := recovered.Outbox()
	if err != nil {
		t.Fatal(err)
	}
	var latest time.Time
	for _, notice := range notices {
		if notice.BlockadeID == blockade.ID && notice.CreatedAt.After(latest) {
			latest = notice.CreatedAt
		}
	}
	if latest.IsZero() {
		t.Fatal("release notices were not recorded")
	}
	checked := 0
	for _, notice := range notices {
		if notice.BlockadeID != blockade.ID || !notice.CreatedAt.Equal(latest) {
			continue
		}
		checked++
		if !strings.Contains(notice.Message, "state=released") {
			t.Fatalf("latest notice for %s used %q", notice.Audience, notice.Message)
		}
	}
	if checked != 4 {
		t.Fatalf("latest notice audience count=%d", checked)
	}
}
