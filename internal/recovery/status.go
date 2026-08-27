package recovery

type Status struct {
	Ready     bool
	Events    int
	Blockades int
	Message   string
}

func NewStatus(r RestartReport) Status {
	return Status{Ready: len(r.Replay.Errors) == 0, Events: r.Replay.Events, Blockades: r.Replay.Blockades, Message: r.Summary()}
}
func (s Status) String() string       { return s.Message }
func (s Status) Operational() bool    { return s.Ready && s.Events >= 0 }
func (s Status) HasEvents() bool      { return s.Events > 0 }
func (s Status) HasBlockades() bool   { return s.Blockades > 0 }
func (s Status) NeedsAttention() bool { return !s.Ready || s.Blockades == 0 }
func (s Status) EventHealth() string {
	if s.Ready {
		return "ok"
	}
	return "degraded"
}
