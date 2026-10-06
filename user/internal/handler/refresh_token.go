package user_transport_grpc

import (
	"context"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) RefreshToken(ctx context.Context, req *userv1.RefreshTokenRequest) (*userv1.RefreshTokenResponse, error) {
	accessToken, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, mapError(err)
	}
	return &userv1.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}
