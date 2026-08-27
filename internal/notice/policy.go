package notice

import "time"

type RetryPolicy struct {
	Limit int
	Base  time.Duration
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	return time.Duration(attempt) * p.Base
}
func (p RetryPolicy) Allowed(attempt int) bool   { return attempt <= p.Limit }
func (p RetryPolicy) Exhausted(attempt int) bool { return !p.Allowed(attempt) }
