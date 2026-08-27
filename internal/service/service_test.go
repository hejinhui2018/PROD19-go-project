package service

import (
	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/domain"
	"testing"
	"time"
)

func TestLifecycle(t *testing.T) {
	d := t.TempDir()
	s, _ := New(d, clock.Fixed{T: time.Now()})
	b := domain.Blockade{ID: "b", Name: "n", Segment: domain.TrackSegment{ID: "s"}, Crew: domain.Crew{ID: "c", Qualified: true}, Start: time.Now(), End: time.Now().Add(time.Hour)}
	if s.Create(b) != nil {
		t.Fatal()
	}
	if s.Transition("b", domain.StatePlanned) != nil {
		t.Fatal()
	}
}
