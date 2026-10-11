package gateway_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	gateway_middleware "github.com/sinon273/url-shortener/gateway/internal/middleware"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := gateway_middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user id", http.StatusUnauthorized)
		return
	}
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		http.Error(w, "missing short code", http.StatusBadRequest)
		return
	}

	req := &urlv1.DeleteURLRequest{
		ShortCode: shortCode,
		UserId:    userID,
	}
	_, err := h.client.URL().DeleteURL(r.Context(), req)
	if err != nil {
		mapGRPCError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
