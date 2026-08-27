package notice

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/domain"
	"strings"
)

type Route struct {
	Audience Audience
	Address  string
	Enabled  bool
}
type Router struct{ routes map[Audience]Route }

func NewRouter(routes []Route) *Router {
	m := map[Audience]Route{}
	for _, r := range routes {
		m[r.Audience] = r
	}
	return &Router{m}
}
func (r *Router) Build(event domain.Event) []domain.ReleaseNotice {
	out := []domain.ReleaseNotice{}
	b := event.Payload
	for a, x := range r.routes {
		if x.Enabled && strings.TrimSpace(x.Address) != "" {
			out = append(out, domain.ReleaseNotice{ID: fmt.Sprintf("%s-%s-%d", event.BlockadeID, a, b.Version), BlockadeID: event.BlockadeID, Audience: string(a), Message: Render(a, event.BlockadeID, string(b.State)), CreatedAt: event.At})
		}
	}
	return out
}
func (r *Router) Has(a Audience) bool { x, ok := r.routes[a]; return ok && x.Enabled }
