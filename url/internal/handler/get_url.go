package url_handler

import (
	"context"

	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

func (h *Handler) GetURL(ctx context.Context, req *urlv1.GetURLRequest) (*urlv1.GetURLResponse, error) {
	link, err := h.service.GetURL(ctx, req.ShortCode)
	if err != nil {
		return nil, mapError(err)
	}

	var expiresAt int64
	if link.ExpiresAt != nil {
		expiresAt = link.ExpiresAt.Unix()
	}

	return &urlv1.GetURLResponse{
		OriginalUrl: link.OriginalURL,
		UserId:      link.UserID,
		CreatedAt:   link.CreatedAt.Unix(),
		ExpiresAt:   expiresAt,
		IsActive:    link.IsActive,
	}, nil
}
