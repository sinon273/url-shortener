package user_transport_grpc

import (
	"context"
	"log/slog"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) Logout(ctx context.Context, req *userv1.LogoutRequest) (*userv1.LogoutResponse, error) {
	err := h.service.Logout(ctx, req.SessionId)
	if err != nil {
		return nil, mapError(ctx, err)
	}
	slog.InfoContext(ctx, "user logged out", "session_id", req.SessionId)
	return &userv1.LogoutResponse{}, nil
}
