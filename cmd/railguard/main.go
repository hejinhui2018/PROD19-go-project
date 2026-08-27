package main

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/domain"
	"github.com/wogo-prod19/railguard/internal/ops"
	"github.com/wogo-prod19/railguard/internal/service"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	if os.Args[1] != "smoke" {
		return
	}
	d, e := os.MkdirTemp("", "railguard-")
	if e != nil {
		panic(e)
	}
	s, e := service.New(d, clock.Real{})
	if e != nil {
		panic(e)
	}
	b := domain.Blockade{ID: ops.ID("blk"), Name: "North work", Segment: domain.TrackSegment{ID: "seg-1", Name: "North"}, Crew: domain.Crew{ID: "crew-1", Name: "A", Qualified: true}, Start: time.Now(), End: time.Now().Add(time.Hour)}
	if e = s.Create(b); e != nil {
		panic(e)
	}
	for _, st := range []domain.State{domain.StatePlanned, domain.StateCrewReady, domain.StateProtected, domain.StateActive, domain.StateReleasePending, domain.StateReleased} {
		if e = s.Transition(b.ID, st); e != nil {
			panic(e)
		}
	}
	s2, e := service.New(d, clock.Real{})
	if e != nil {
		panic(e)
	}
	x, e := s2.Get(b.ID)
	if e != nil || x.State != domain.StateReleased {
		panic("recovery failed")
	}
	if _, e = s.Plan(b.ID, "dispatcher"); e != nil {
		panic(e)
	}
	if _, e = s.RecoveryReport(); e != nil {
		panic(e)
	}
	if d, e := s.DrainNotices(); e != nil || !d.Complete() {
		panic("outbox drain failed")
	}
	fmt.Println("railguard smoke: lifecycle and recovery succeeded")
}
