package analytics_domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrPermissionDenied = errors.New("permission denied")
	ErrUnavailable      = errors.New("unavailable")
)
