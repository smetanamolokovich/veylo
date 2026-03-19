package finding

import "errors"

var (
	ErrNotFound            = errors.New("finding: not found")
	ErrInvalidSeverity     = errors.New("finding: invalid severity")
	ErrInvalidRepairMethod = errors.New("finding: invalid repair method")
)
