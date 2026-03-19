package asset

import "errors"

var (
	ErrNotFound      = errors.New("asset: not found")
	ErrAlreadyExists = errors.New("asset: already exists")
	ErrInvalidType   = errors.New("asset: invalid type")
)
