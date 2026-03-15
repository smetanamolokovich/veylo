package inspection

import "errors"

var (
	ErrNotFound          = errors.New("inspection: not found")
	ErrInvalidTransition = errors.New("inspection: invalid status transition")
	ErrAlreadyCompleted  = errors.New("inspection: already completed")
)
