package gateway_handler

import (
	"encoding/json"
	"net/http"

	gateway_middleware "github.com/sinon273/url-shortener/gateway/internal/middleware"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

type createURLRequest struct {
	OriginalURL string `json:"original_url"`
	ExpiresIn   *int64 `json:"expires_in,omitempty"`
}

func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := gateway_middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user id", http.StatusUnauthorized)
		return
	}
	var body createURLRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req := urlv1.CreateShortURLRequest{
		OriginalUrl: body.OriginalURL,
		UserId:      userID,
		ExpiresIn:   body.ExpiresIn,
	}

	resp, err := h.client.URL().CreateShortURL(r.Context(), &req)
	if err != nil {
		mapGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)

}
