package planning

import (
	"fmt"
	"time"
)

type WorkWindow struct {
	Start, End   time.Time
	Possessions  []string
	WeatherClear bool
}

func (w WorkWindow) Validate() error {
	if w.Start.IsZero() || w.End.IsZero() || !w.Start.Before(w.End) {
		return fmt.Errorf("invalid work window")
	}
	if len(w.Possessions) == 0 {
		return fmt.Errorf("track possession required")
	}
	return nil
}
func (w WorkWindow) Contains(t time.Time) bool { return !t.Before(w.Start) && t.Before(w.End) }
func (w WorkWindow) Duration() time.Duration   { return w.End.Sub(w.Start) }
