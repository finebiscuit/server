package payload

import "errors"

var (
	ErrInvalidScheme = errors.New("invalid scheme")
	ErrEmptyVersion  = errors.New("missing version")
)
