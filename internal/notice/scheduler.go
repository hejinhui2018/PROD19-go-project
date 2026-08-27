package notice

import "time"

type Schedule struct {
	At        time.Time
	NoticeIDs []string
}

func (s Schedule) Due(now time.Time) bool { return !s.At.IsZero() && !now.Before(s.At) }
func (s Schedule) Remaining(now time.Time) time.Duration {
	if s.Due(now) {
		return 0
	}
	return s.At.Sub(now)
}
