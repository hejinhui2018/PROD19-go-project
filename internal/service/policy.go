package service

import "github.com/wogo-prod19/railguard/internal/domain"

func Terminal(s domain.State) bool {
	return s == domain.StateReleased || s == domain.StateCancelled || s == domain.StateIncident
}
