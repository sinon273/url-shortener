package domain

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrNotFound           = errors.New("not found")
	ErrUnauthenticated    = errors.New("unauthenticated")
)
