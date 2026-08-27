package store

import (
	"github.com/wogo-prod19/railguard/internal/domain"
	"testing"
)

func TestOutbox(t *testing.T) {
	o := NewOutbox(t.TempDir())
	_ = o.Add(domain.ReleaseNotice{ID: "n"})
	x, _ := o.Drain()
	if len(x) != 1 {
		t.Fatal()
	}
}
