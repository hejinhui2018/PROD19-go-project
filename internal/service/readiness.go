package service

import "github.com/wogo-prod19/railguard/internal/domain"

func ReadyForProtection(b domain.Blockade) bool {
	return b.Crew.Qualified && b.Approval.Approved && b.Locks.Applied && b.Equipment.Complete
}
