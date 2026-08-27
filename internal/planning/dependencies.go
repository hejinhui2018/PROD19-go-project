package planning

import "github.com/wogo-prod19/railguard/internal/domain"

func CheckAdjacent(segments []domain.TrackSegment, links []domain.SegmentLink) []DependencyIssue {
	if len(segments) < 2 {
		return nil
	}
	have := map[string]bool{}
	for _, s := range segments {
		have[s.ID] = true
	}
	issues := []DependencyIssue{}
	for _, l := range links {
		if !have[l.From] || !have[l.To] {
			issues = append(issues, DependencyIssue{l.From, l.To, "segment missing from draft"})
		}
		if l.From == l.To {
			issues = append(issues, DependencyIssue{l.From, l.To, "self dependency"})
		}
		if l.Distance <= 0 {
			issues = append(issues, DependencyIssue{l.From, l.To, "distance must be positive"})
		}
	}
	return issues
}
func BuildLinks(segments []domain.TrackSegment) []domain.SegmentLink {
	links := make([]domain.SegmentLink, 0, len(segments)-1)
	for i := 0; i+1 < len(segments); i++ {
		links = append(links, domain.SegmentLink{From: segments[i].ID, To: segments[i+1].ID, Distance: i + 1})
	}
	return links
}
