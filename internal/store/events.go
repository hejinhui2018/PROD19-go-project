package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/wogo-prod19/railguard/internal/domain"
	"os"
	"path/filepath"
)

type EventStore struct {
	events, snapshot string
	next             int64
}

func Open(dir string) (*EventStore, error) {
	if e := os.MkdirAll(dir, 0755); e != nil {
		return nil, e
	}
	s := &EventStore{events: filepath.Join(dir, "events.jsonl"), snapshot: filepath.Join(dir, "snapshot.json")}
	if e := s.scan(); e != nil {
		return nil, e
	}
	return s, nil
}
func (s *EventStore) scan() error {
	f, e := os.Open(s.events)
	if os.IsNotExist(e) {
		return nil
	}
	if e != nil {
		return e
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var n int64
	for sc.Scan() {
		var ev domain.Event
		if json.Unmarshal(sc.Bytes(), &ev) != nil || ev.Number <= n {
			return fmt.Errorf("%w: event %d", domain.ErrCorrupt, n+1)
		}
		n = ev.Number
	}
	s.next = n
	return sc.Err()
}
func (s *EventStore) Append(ev domain.Event) error {
	s.next++
	ev.Number = s.next
	b, e := json.Marshal(ev)
	if e != nil {
		return e
	}
	f, e := os.OpenFile(s.events, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if e != nil {
		return e
	}
	defer f.Close()
	_, e = f.Write(append(b, '\n'))
	return e
}
func (s *EventStore) Load() ([]domain.Event, error) {
	f, e := os.Open(s.events)
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	defer f.Close()
	var out []domain.Event
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var ev domain.Event
		if json.Unmarshal(sc.Bytes(), &ev) != nil {
			return nil, domain.ErrCorrupt
		}
		out = append(out, ev)
	}
	return out, sc.Err()
}
func (s *EventStore) SaveSnapshot(b domain.Blockade) error {
	x, e := json.Marshal(b)
	if e != nil {
		return e
	}
	return os.WriteFile(s.snapshot, x, 0644)
}
func (s *EventStore) LoadSnapshot() (domain.Blockade, error) {
	x, e := os.ReadFile(s.snapshot)
	if e != nil {
		return domain.Blockade{}, e
	}
	var b domain.Blockade
	e = json.Unmarshal(x, &b)
	return b, e
}
