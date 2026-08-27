package safety

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/domain"
	"time"
)

type LockCheck struct {
	LockID, Inspector string
	VerifiedAt        time.Time
	Evidence          []string
}

func VerifyLockPack(pack domain.SafetyLockPack, check LockCheck) error {
	if pack.ID == "" || check.LockID != pack.ID {
		return fmt.Errorf("lock identity mismatch")
	}
	if !pack.Applied {
		return fmt.Errorf("lock pack not applied")
	}
	if check.Inspector == "" || len(check.Evidence) == 0 {
		return fmt.Errorf("lock evidence incomplete")
	}
	return nil
}
