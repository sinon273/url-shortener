package url_handler

import (
	"context"

	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

func (h *Handler) CreateShortURL(ctx context.Context, req *urlv1.CreateShortURLRequest) (*urlv1.CreateShortURLResponse, error) {
	shortCode, shortURL, err := h.service.CreateShortURL(ctx, req.OriginalUrl, req.UserId, req.ExpiresIn)
	if err != nil {
		return nil, mapError(err)
	}

	return &urlv1.CreateShortURLResponse{
		ShortCode: shortCode,
		ShortUrl:  shortURL,
	}, nil
}
