package workflow

import "errors"

var (
	ErrNotFound                = errors.New("workflow: not found")
	ErrNoInitialStatus         = errors.New("workflow: no initial status defined")
	ErrInitialStatusAlreadySet = errors.New("workflow: initial status already set")
	ErrDuplicateStatus         = errors.New("workflow: status already exists")
	ErrStatusNotFound          = errors.New("workflow: status not found")
)
