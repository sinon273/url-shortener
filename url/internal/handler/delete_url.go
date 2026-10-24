package url_handler

import (
	"context"
	"log/slog"

	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

func (h *Handler) DeleteURL(ctx context.Context, req *urlv1.DeleteURLRequest) (*urlv1.DeleteURLResponse, error) {
	err := h.service.DeleteURL(ctx, req.ShortCode, req.UserId)
	if err != nil {
		return nil, mapError(ctx, err)
	}

	slog.InfoContext(ctx, "deleted url", "short_code", req.ShortCode)
	return &urlv1.DeleteURLResponse{}, nil
}
