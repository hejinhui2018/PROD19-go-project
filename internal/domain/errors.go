package domain

import "errors"

var (
	ErrNotFound = errors.New("blockade not found")
	ErrInvalid  = errors.New("invalid blockade")
	ErrConflict = errors.New("invalid state transition")
	ErrCorrupt  = errors.New("corrupt event record")
)
