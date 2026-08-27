package notice

import "github.com/wogo-prod19/railguard/internal/domain"

type AudienceStats struct{ Sent, Failed int }
type DeliveryAnalytics struct{ ByAudience map[Audience]AudienceStats }

func (a *DeliveryAnalytics) Add(n domain.ReleaseNotice, ok bool) {
	if a.ByAudience == nil {
		a.ByAudience = map[Audience]AudienceStats{}
	}
	x := a.ByAudience[Audience(n.Audience)]
	if ok {
		x.Sent++
	} else {
		x.Failed++
	}
	a.ByAudience[Audience(n.Audience)] = x
}
func (a DeliveryAnalytics) Total() int {
	n := 0
	for _, x := range a.ByAudience {
		n += x.Sent + x.Failed
	}
	return n
}
func (a DeliveryAnalytics) Failed() int {
	n := 0
	for _, x := range a.ByAudience {
		n += x.Failed
	}
	return n
}
func (a DeliveryAnalytics) Healthy() bool { return a.Failed() == 0 }
