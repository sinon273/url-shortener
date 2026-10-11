package gateway_handler

import (
	"encoding/json"
	"net/http"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req userv1.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.client.User().Login(ctx, &req)
	if err != nil {
		mapGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
