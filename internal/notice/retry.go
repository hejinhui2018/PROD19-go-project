package notice

import "time"

type RetryBook struct {
	Attempts  int
	LastError string
	NextAt    time.Time
	Delivered bool
}

func (b *RetryBook) Fail(err string, now time.Time) {
	b.Attempts++
	b.LastError = err
	b.NextAt = now.Add(time.Duration(b.Attempts) * time.Minute)
}
func (b *RetryBook) Success() { b.Delivered = true; b.LastError = "" }
func (b RetryBook) Retryable(now time.Time) bool {
	return !b.Delivered && b.Attempts < 8 && !now.Before(b.NextAt)
}
