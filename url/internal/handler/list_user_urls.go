package url_handler

import (
	"context"

	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

func (h *Handler) ListUserURLs(ctx context.Context, req *urlv1.ListUserURLsRequest) (*urlv1.ListUserURLsResponse, error) {
	links, nextCursor, err := h.service.ListUserURLs(ctx, req.UserId, int(req.Limit), req.Cursor)
	if err != nil {
		return nil, mapError(err)
	}

	items := make([]*urlv1.URLItem, len(links))
	for i, link := range links {
		var expiresAt int64
		if link.ExpiresAt != nil {
			expiresAt = link.ExpiresAt.Unix()
		}
		items[i] = &urlv1.URLItem{
			ShortCode:   link.ShortCode,
			OriginalUrl: link.OriginalURL,
			CreatedAt:   link.CreatedAt.Unix(),
			ExpiresAt:   expiresAt,
			IsActive:    link.IsActive,
		}
	}

	return &urlv1.ListUserURLsResponse{
		Urls:       items,
		NextCursor: nextCursor,
	}, nil
}
