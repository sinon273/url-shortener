package gateway_handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	gateway_middleware "github.com/sinon273/url-shortener/gateway/internal/middleware"
	analyticsv1 "github.com/sinon273/url-shortener/shared/pkg/proto/analytics/v1"
)

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
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

	var days int32
	dayStr := r.URL.Query().Get("days")
	if dayStr != "" {
		parsed, err := strconv.Atoi(dayStr)
		if err != nil {
			http.Error(w, "invalid days", http.StatusBadRequest)
			return
		}
		days = int32(parsed)
	}

	resp, err := h.client.Analytics().GetAnalytics(r.Context(), &analyticsv1.GetAnalyticsRequest{
		ShortCode: shortCode,
		UserId:    userID,
		Days:      days,
	})
	if err != nil {
		mapGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
