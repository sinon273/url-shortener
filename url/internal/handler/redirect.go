package url_handler

import (
	"context"

	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

func (h *Handler) Redirect(ctx context.Context, req *urlv1.RedirectRequest) (*urlv1.RedirectResponse, error) {
	originalURL, err := h.service.Redirect(ctx, req.ShortCode)
	if err != nil {
		return nil, mapError(err)
	}
	return &urlv1.RedirectResponse{
		OriginalUrl: originalURL,
	}, nil
}
