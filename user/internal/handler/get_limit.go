package user_transport_grpc

import (
	"context"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) GetLimit(ctx context.Context, req *userv1.GetLimitRequest) (*userv1.GetLimitResponse, error) {
	limit, err := h.service.GetLimit(ctx, req.UserId)
	if err != nil {
		return nil, mapError(err)
	}

	return &userv1.GetLimitResponse{
		LinksLimit: int32(limit),
	}, nil
}
