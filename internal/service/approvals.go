package service

import (
	"github.com/wogo-prod19/railguard/internal/domain"
	"time"
)

func Approve(b domain.Blockade, dispatcher string, now time.Time) domain.Blockade {
	b.Approval = domain.DispatcherApproval{ID: b.ID + "-approval", Dispatcher: dispatcher, Approved: true, At: now}
	return b
}
