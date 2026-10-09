package url_handler

import (
	"errors"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, url_domain.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, url_domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, url_domain.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, url_domain.ErrResourceExhausted):
		return status.Error(codes.ResourceExhausted, err.Error())
	case errors.Is(err, url_domain.ErrUnavailable):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
