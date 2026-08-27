package recovery

import (
	"github.com/wogo-prod19/railguard/internal/domain"
	"sort"
)

type RetentionPolicy struct {
	KeepVersions int
	KeepDuration int64
}

func (p RetentionPolicy) Select(events []domain.Event) []domain.Event {
	out := append([]domain.Event(nil), events...)
	sort.Slice(out, func(i, j int) bool { return out[i].Number > out[j].Number })
	if p.KeepVersions > 0 && len(out) > p.KeepVersions {
		out = out[:p.KeepVersions]
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out
}
func (p RetentionPolicy) ShouldKeep(version int64, current int64) bool {
	return p.KeepVersions <= 0 || current-version < int64(p.KeepVersions)
}

type Integrity struct {
	Checked, Valid int
	Failures       []string
}

func (i *Integrity) Add(ok bool, detail string) {
	i.Checked++
	if ok {
		i.Valid++
	} else {
		i.Failures = append(i.Failures, detail)
	}
}
func (i Integrity) Healthy() bool { return i.Checked == i.Valid }
