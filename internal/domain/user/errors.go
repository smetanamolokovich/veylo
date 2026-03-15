package user

import "errors"

var (
	ErrNotFound           = errors.New("user: not found")
	ErrAlreadyExists      = errors.New("user: already exists")
	ErrInvalidRole        = errors.New("user: invalid role")
	ErrBlocked            = errors.New("user: blocked")
	ErrInvalidCredentials = errors.New("user: invalid credentials")
)
