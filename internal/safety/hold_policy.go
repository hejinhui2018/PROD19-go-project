package safety

import (
	"fmt"
	"time"
)

type HoldPolicy struct {
	MaxMinutes         int
	RequiresSignalDesk bool
	ReleaseWindow      time.Duration
}

func (p HoldPolicy) Validate(h Hold, now time.Time) error {
	if !h.Active {
		return fmt.Errorf("hold inactive")
	}
	if p.MaxMinutes > 0 && now.Sub(time.Now()) > time.Duration(p.MaxMinutes)*time.Minute {
		return fmt.Errorf("hold exceeded policy")
	}
	return nil
}
func CanRelease(h Hold, board ProtectionBoard, incidentOpen bool) bool {
	return h.Active && board.Active && !incidentOpen
}
