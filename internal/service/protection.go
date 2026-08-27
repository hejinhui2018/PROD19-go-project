package service

import (
	"github.com/wogo-prod19/railguard/internal/domain"
	"time"
)

func Protect(b domain.Blockade, lockID string, now time.Time) domain.Blockade {
	b.Locks = domain.SafetyLockPack{ID: lockID, Applied: true, AppliedAt: now}
	return b
}
