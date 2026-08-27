package notice

import "time"

type Receipt struct {
	NoticeID, Address string
	Accepted          bool
	At                time.Time
	Detail            string
}

func (r Receipt) Valid() bool       { return r.NoticeID != "" && r.Address != "" && !r.At.IsZero() }
func (r Receipt) RetryNeeded() bool { return !r.Accepted }
