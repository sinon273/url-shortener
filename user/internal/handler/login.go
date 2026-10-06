package user_transport_grpc

import (
	"context"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	accessToken, refreshToken, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, mapError(err)
	}

	return &userv1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
