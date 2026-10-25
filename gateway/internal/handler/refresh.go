package gateway_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req userv1.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.client.User().RefreshToken(ctx, &req)
	if err != nil {
		mapGRPCError(w, r, err)
		return
	}

	slog.InfoContext(r.Context(), "user refreshed")
	writeJSON(w, http.StatusOK, resp)
}
