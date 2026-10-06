package user_transport_grpc

import (
	"context"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) Logout(ctx context.Context, req *userv1.LogoutRequest) (*userv1.LogoutResponse, error) {
	err := h.service.Logout(ctx, req.SessionId)
	if err != nil {
		return nil, mapError(err)
	}
	return &userv1.LogoutResponse{}, nil
}
