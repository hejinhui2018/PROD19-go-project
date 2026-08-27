package domain

func ValidState(s State) bool {
	switch s {
	case StateDraft, StatePlanned, StateCrewReady, StateProtected, StateActive, StateReleasePending, StateReleased, StateCancelled, StateIncident:
		return true
	}
	return false
}
