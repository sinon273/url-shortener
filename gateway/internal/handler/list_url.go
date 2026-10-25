package gateway_handler

import (
	"net/http"
	"strconv"
	"time"

	gateway_middleware "github.com/sinon273/url-shortener/gateway/internal/middleware"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

type urlItemResponse struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
	ExpiresAt   string `json:"expires_at"`
	IsActive    bool   `json:"is_active"`
}

type listURLsResponse struct {
	URLs       []urlItemResponse `json:"urls"`
	NextCursor string            `json:"next_cursor"`
}

func (h *Handler) ListURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := gateway_middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user id", http.StatusUnauthorized)
		return
	}
	limitStr := r.URL.Query().Get("limit")
	var limit int32 = 0
	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = int32(parsed)
	}

	cursor := r.URL.Query().Get("cursor")

	resp, err := h.client.URL().ListUserURLs(r.Context(), &urlv1.ListUserURLsRequest{
		UserId: userID,
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		mapGRPCError(w, r, err)
		return
	}

	converted := make([]urlItemResponse, len(resp.Urls))
	for i, item := range resp.Urls {

		expiresAt := ""
		if item.ExpiresAt != 0 {
			expiresAt = time.Unix(item.ExpiresAt, 0).UTC().Format(time.RFC3339)
		}

		converted[i] = urlItemResponse{
			ShortCode:   item.ShortCode,
			OriginalURL: item.OriginalUrl,
			CreatedAt:   time.Unix(item.CreatedAt, 0).UTC().Format(time.RFC3339),
			ExpiresAt:   expiresAt,
			IsActive:    item.IsActive,
		}
	}
	list := listURLsResponse{
		URLs:       converted,
		NextCursor: resp.NextCursor,
	}

	writeJSON(w, http.StatusOK, list)
}
