package service

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/audit"
	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/domain"
	"github.com/wogo-prod19/railguard/internal/notice"
	"github.com/wogo-prod19/railguard/internal/planning"
	"github.com/wogo-prod19/railguard/internal/recovery"
	"github.com/wogo-prod19/railguard/internal/safety"
	"github.com/wogo-prod19/railguard/internal/store"
	"sync"
)

type Service struct {
	mu     sync.Mutex
	blocks map[string]domain.Blockade
	events *store.EventStore
	outbox *store.Outbox
	clk    clock.Clock
	audit  *audit.Ledger
	router *notice.Router
}

func New(dir string, c clock.Clock) (*Service, error) {
	es, e := store.Open(dir)
	if e != nil {
		return nil, e
	}
	s := &Service{blocks: map[string]domain.Blockade{}, events: es, outbox: store.NewOutbox(dir), clk: c, audit: &audit.Ledger{}, router: notice.NewRouter([]notice.Route{{Audience: notice.Dispatcher, Address: "dispatcher", Enabled: true}, {Audience: notice.Crew, Address: "crew", Enabled: true}, {Audience: notice.SignalDesk, Address: "signal", Enabled: true}, {Audience: notice.TrackAuthority, Address: "authority", Enabled: true}})}
	evs, e := es.Load()
	if e != nil {
		return nil, e
	}
	if _, e = recovery.VerifyReplay(evs); e != nil && len(evs) > 0 {
		return nil, e
	}
	for _, ev := range evs {
		s.blocks[ev.BlockadeID] = ev.Payload
	}
	return s, nil
}
func (s *Service) Create(b domain.Blockade) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e := domain.ValidateBlockade(b); e != nil {
		return e
	}
	if _, ok := s.blocks[b.ID]; ok {
		return fmt.Errorf("%w: duplicate", domain.ErrConflict)
	}
	b.State = domain.StateDraft
	b.UpdatedAt = s.clk.Now()
	return s.persist("created", b)
}
func (s *Service) Transition(id string, to domain.State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.blocks[id]
	if !ok {
		return domain.ErrNotFound
	}
	if !allowed(b.State, to) {
		return fmt.Errorf("%w: %s to %s", domain.ErrConflict, b.State, to)
	}
	b.State = to
	b.Version++
	b.UpdatedAt = s.clk.Now()
	notification := domain.Event{Type: domain.EventTransition, BlockadeID: b.ID, Payload: b, At: b.UpdatedAt}
	if e := s.persist("transition", b); e != nil {
		return e
	}
	_ = s.audit.Append(audit.New(fmt.Sprintf("%s-%d", id, b.Version), "dispatcher", "transition", id, string(b.State), string(to), "operational state transition", b.UpdatedAt))
	for _, n := range s.router.Build(notification) {
		if err := s.outbox.Add(n); err != nil {
			return err
		}
	}
	return nil
}
func allowed(a, b domain.State) bool {
	m := map[domain.State]domain.State{domain.StateDraft: domain.StatePlanned, domain.StatePlanned: domain.StateCrewReady, domain.StateCrewReady: domain.StateProtected, domain.StateProtected: domain.StateActive, domain.StateActive: domain.StateReleasePending, domain.StateReleasePending: domain.StateReleased}
	return m[a] == b || b == domain.StateCancelled || b == domain.StateIncident
}
func (s *Service) persist(t string, b domain.Blockade) error {
	if e := s.events.Append(domain.Event{Type: t, BlockadeID: b.ID, Payload: b, At: b.UpdatedAt}); e != nil {
		return e
	}
	s.blocks[b.ID] = b
	return s.events.SaveSnapshot(b)
}
func (s *Service) Get(id string) (domain.Blockade, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.blocks[id]
	if !ok {
		return domain.Blockade{}, domain.ErrNotFound
	}
	return b, nil
}
func (s *Service) Outbox() ([]domain.ReleaseNotice, error) { return s.outbox.List() }
func (s *Service) AddObservation(id string, o domain.FieldObservation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.blocks[id]
	if !ok {
		return domain.ErrNotFound
	}
	b.Observations = append(b.Observations, o)
	b.UpdatedAt = s.clk.Now()
	return s.persist("observation", b)
}

func (s *Service) Plan(id, author string) (planning.Draft, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.blocks[id]
	if !ok {
		return planning.Draft{}, domain.ErrNotFound
	}
	d := planning.NewDraft("plan-"+id, id, author, []domain.TrackSegment{b.Segment}, s.clk.Now())
	d.Dependencies = planning.BuildLinks(d.Segments)
	d.Crew = planning.CrewPlan{CrewID: b.Crew.ID, Qualified: b.Crew.Qualified, Shift: "day", Briefed: true}
	d.Equipment = planning.EquipmentPlan{Required: []string{"lock-kit"}, Available: []string{"lock-kit"}, Inspected: true, Custodian: b.Crew.ID}
	if err := d.Validate(); err != nil {
		return d, err
	}
	return d, nil
}
func (s *Service) VerifySafety(id string) error {
	b, e := s.Get(id)
	if e != nil {
		return e
	}
	if !ReadyForProtection(b) {
		return fmt.Errorf("%w: safety readiness", domain.ErrConflict)
	}
	_, e = safety.IssueBoard("board-"+id, id, "dispatcher", []safety.BoardEntry{{Signal: "SIG-1", Aspect: "stop", Route: b.Segment.ID, Locked: true}}, s.clk.Now())
	return e
}
func (s *Service) DrainNotices() (notice.DrainResult, error) {
	ns, e := s.outbox.Drain()
	if e != nil {
		return notice.DrainResult{}, e
	}
	r := notice.NewDrain(s.clk.Now())
	r.Seen = len(ns)
	r.Delivered = len(ns)
	return r, nil
}
func (s *Service) RecoveryReport() (recovery.RestartReport, error) {
	ev, e := s.events.Load()
	if e != nil {
		return recovery.RestartReport{}, e
	}
	rr, e := recovery.VerifyReplay(ev)
	return recovery.RestartReport{StartedAt: s.clk.Now(), CompletedAt: s.clk.Now(), Replay: rr}, e
}
