package user

import "errors"

var (
	ErrNotFound           = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrAlreadyExists      = errors.New("user already exists")
	ErrEmailAlreadyTaken  = errors.New("email is already taken")
)
