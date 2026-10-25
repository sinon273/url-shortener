package user_transport_grpc

import (
	"context"
	"log/slog"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) ValidateSession(ctx context.Context, req *userv1.ValidateSessionRequest) (*userv1.ValidateSessionResponse, error) {
	userID, sessionID, err := h.service.ValidateSession(ctx, req.AccessToken)
	if err != nil {
		return nil, mapError(ctx, err)
	}

	slog.InfoContext(ctx, "session validated", "user_id", userID)
	return &userv1.ValidateSessionResponse{
		UserId:    userID,
		SessionId: sessionID,
	}, nil
}
