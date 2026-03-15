package refreshtoken

import "errors"

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrExpiredRefreshToken = errors.New("expired refresh token")
	ErrNotFound            = errors.New("refresh token: not found")
)
