package notice

import "github.com/wogo-prod19/railguard/internal/domain"

type NoticeBook struct {
	Pending map[string]domain.ReleaseNotice
	History []Delivery
}

func (b *NoticeBook) Put(n domain.ReleaseNotice) {
	if b.Pending == nil {
		b.Pending = map[string]domain.ReleaseNotice{}
	}
	b.Pending[n.ID] = n
}
func (b *NoticeBook) Mark(d Delivery) {
	delete(b.Pending, d.NoticeID)
	b.History = append(b.History, d)
}
func (b NoticeBook) PendingCount() int { return len(b.Pending) }
func (b NoticeBook) SuccessRate() float64 {
	if len(b.History) == 0 {
		return 1
	}
	ok := 0
	for _, d := range b.History {
		if d.Successful() {
			ok++
		}
	}
	return float64(ok) / float64(len(b.History))
}
