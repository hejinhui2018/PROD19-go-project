package store

import (
	"github.com/wogo-prod19/railguard/internal/domain"
)

type Repository struct{ Events *EventStore }

func NewRepository(e *EventStore) *Repository         { return &Repository{Events: e} }
func (r *Repository) Replay() ([]domain.Event, error) { return r.Events.Load() }
