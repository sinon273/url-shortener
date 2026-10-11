package gateway_handler

import (
	"net/http"

	gateway_middleware "github.com/sinon273/url-shortener/gateway/internal/middleware"
	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := gateway_middleware.SessionIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	req := userv1.LogoutRequest{SessionId: sessionID}

	resp, err := h.client.User().Logout(r.Context(), &req)
	if err != nil {
		mapGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
