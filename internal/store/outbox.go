package store

import (
	"encoding/json"
	"github.com/wogo-prod19/railguard/internal/domain"
	"os"
	"path/filepath"
)

type Outbox struct{ path string }

func NewOutbox(dir string) *Outbox { return &Outbox{filepath.Join(dir, "outbox.json")} }
func (o *Outbox) Add(n domain.ReleaseNotice) error {
	all, _ := o.List()
	all = append(all, n)
	b, e := json.Marshal(all)
	if e != nil {
		return e
	}
	return os.WriteFile(o.path, b, 0644)
}
func (o *Outbox) List() ([]domain.ReleaseNotice, error) {
	b, e := os.ReadFile(o.path)
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	var n []domain.ReleaseNotice
	e = json.Unmarshal(b, &n)
	return n, e
}
func (o *Outbox) Drain() ([]domain.ReleaseNotice, error) {
	n, e := o.List()
	if e != nil {
		return nil, e
	}
	if e = os.WriteFile(o.path, []byte("[]"), 0644); e != nil {
		return nil, e
	}
	return n, nil
}
