package organization

import "errors"

var (
	ErrNotFound      = errors.New("organization: not found")
	ErrAlreadyExists = errors.New("organization: already exists")
)
