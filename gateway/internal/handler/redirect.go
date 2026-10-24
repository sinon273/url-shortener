package gateway_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		http.Error(w, "missing short code", http.StatusBadRequest)
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	userAgent := r.Header.Get("User-Agent")
	referer := r.Header.Get("Referer")

	req := &urlv1.RedirectRequest{
		ShortCode: shortCode,
		Ip:        ip,
		UserAgent: userAgent,
		Referer:   referer,
	}

	resp, err := h.client.URL().Redirect(r.Context(), req)
	if err != nil {
		mapGRPCError(w, r, err)
		return
	}

	w.Header().Set("Location", resp.OriginalUrl)
	w.WriteHeader(http.StatusFound)
}
