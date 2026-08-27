package notice

import "fmt"

type Audience string

const (
	Dispatcher     Audience = "dispatcher"
	Crew           Audience = "crew"
	SignalDesk     Audience = "signal-desk"
	TrackAuthority Audience = "track-authority"
)

func Render(a Audience, blockade, state string) string {
	return fmt.Sprintf("[%s] blockade %s requires attention: state=%s", a, blockade, state)
}
