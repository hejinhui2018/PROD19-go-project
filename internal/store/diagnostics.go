package store

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/domain"
)

type Diagnostics struct {
	Count int64
	Last  int64
}

func (s *EventStore) Diagnostics() (Diagnostics, error) {
	e, err := s.Load()
	if err != nil {
		return Diagnostics{}, err
	}
	d := Diagnostics{Count: int64(len(e))}
	if len(e) > 0 {
		d.Last = e[len(e)-1].Number
	}
	if d.Count > 0 && d.Last != d.Count {
		return d, fmt.Errorf("%w: non-monotonic", domain.ErrCorrupt)
	}
	return d, nil
}
