package url_domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrResourceExhausted = errors.New("resource exhausted")
	ErrUnavailable       = errors.New("unavailable")
	ErrAlreadyExists     = errors.New("already exists")
)
