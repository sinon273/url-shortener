package user_transport_grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/sinon273/url-shortener/user/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(ctx context.Context, err error) error {
	slog.ErrorContext(ctx, "request failed", "err", err)
	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
