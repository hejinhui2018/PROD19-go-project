package service

import (
	"github.com/wogo-prod19/railguard/internal/domain"
	"time"
)

func Handoff(b domain.Blockade, id, custodian string, now time.Time) domain.Blockade {
	b.Equipment = domain.EquipmentHandoff{ID: id, Custodian: custodian, Complete: true, At: now}
	return b
}
