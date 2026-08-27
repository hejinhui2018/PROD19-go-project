package notice

import "github.com/wogo-prod19/railguard/internal/domain"

func GroupByAudience(ns []domain.ReleaseNotice) map[Audience][]domain.ReleaseNotice {
	out := map[Audience][]domain.ReleaseNotice{}
	for _, n := range ns {
		a := Audience(n.Audience)
		out[a] = append(out[a], n)
	}
	return out
}
