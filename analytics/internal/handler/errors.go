package analytics_handler

import (
	"context"
	"errors"
	"log/slog"

	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(ctx context.Context, err error) error {
	slog.ErrorContext(ctx, "request failed", "err", err)

	switch {
	case errors.Is(err, analytics_domain.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, analytics_domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, analytics_domain.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, analytics_domain.ErrUnavailable):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
