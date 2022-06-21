package balance

import "errors"

var (
	ErrNotFound       = errors.New("balance not found")
	ErrInvalidPayload = errors.New("invalid payload")
)
