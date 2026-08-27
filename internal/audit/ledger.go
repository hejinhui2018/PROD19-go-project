package audit

import (
	"fmt"
	"sync"
)

type Ledger struct {
	mu      sync.Mutex
	records []Record
}

func (l *Ledger) Append(r Record) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if r.ID == "" || r.Hash == "" {
		return fmt.Errorf("invalid audit record")
	}
	for _, x := range l.records {
		if x.ID == r.ID {
			return fmt.Errorf("duplicate audit record")
		}
	}
	l.records = append(l.records, r)
	return nil
}
func (l *Ledger) List() []Record {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]Record(nil), l.records...)
}
