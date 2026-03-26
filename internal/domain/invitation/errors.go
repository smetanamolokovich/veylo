package invitation

import "errors"

var (
	ErrNotFound    = errors.New("invitation: not found")
	ErrAlreadyUsed = errors.New("invitation: already used")
	ErrExpired     = errors.New("invitation: expired")
	ErrDuplicate   = errors.New("invitation: pending invitation already exists for this email")
)
