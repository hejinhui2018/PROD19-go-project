package store

import (
	"github.com/wogo-prod19/railguard/internal/domain"
	"testing"
)

func TestEvents(t *testing.T) {
	s, _ := Open(t.TempDir())
	if s.Append(domain.Event{Type: "x", BlockadeID: "b"}) != nil {
		t.Fatal()
	}
	d, e := s.Diagnostics()
	if e != nil || d.Count != 1 {
		t.Fatal(d, e)
	}
}
