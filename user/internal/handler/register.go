package user_transport_grpc

import (
	"context"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) Register(
	ctx context.Context,
	req *userv1.RegisterRequest,
) (*userv1.RegisterResponse, error) {
	userID, err := h.service.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, mapError(err)
	}

	return &userv1.RegisterResponse{
		UserId: userID,
	}, nil
}
