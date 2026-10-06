package user_transport_grpc

import (
	"context"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

type Handler struct {
	userv1.UnimplementedUserServiceServer
	service Service
}

func New(service Service) *Handler {
	return &Handler{service: service}
}

type Service interface {
	Register(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	Logout(ctx context.Context, sessionID string) error
	ValidateSession(ctx context.Context, accessToken string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetLimit(ctx context.Context, userID string) (int, error)
}
