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
func (r *Router) Build(b domain.Blockade) []domain.ReleaseNotice {
	out := []domain.ReleaseNotice{}
	for a, x := range r.routes {
		if x.Enabled && strings.TrimSpace(x.Address) != "" {
			out = append(out, domain.ReleaseNotice{ID: fmt.Sprintf("%s-%s-%d", b.ID, a, b.Version), BlockadeID: b.ID, Audience: string(a), Message: Render(a, b.ID, string(b.State)), CreatedAt: b.UpdatedAt})
		}
	}
	return out
}
func (r *Router) Has(a Audience) bool { x, ok := r.routes[a]; return ok && x.Enabled }
