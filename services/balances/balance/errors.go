package balance

import "errors"

var (
	ErrNotFound           = errors.New("balance not found")
	ErrEntryNotFound      = errors.New("entry not found")
	ErrEntryAlreadyExists = errors.New("entry already exists")
	ErrVersionMismatch    = errors.New("version mismatch")
	ErrInvalidPayload     = errors.New("invalid payload")
)
