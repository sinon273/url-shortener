package gateway_handler

import (
	"github.com/go-chi/chi/v5"
	gateway_middleware "github.com/sinon273/url-shortener/gateway/internal/middleware"
)

func (h *Handler) Router(authMW *gateway_middleware.Auth) *chi.Mux {
	r := chi.NewRouter()

	r.Use(gateway_middleware.MetricsMiddleware)

	r.Post("/api/v1/auth/register", h.Register)
	r.Post("/api/v1/auth/login", h.Login)
	r.Post("/api/v1/auth/refresh", h.Refresh)

	r.Group(func(r chi.Router) {
		r.Use(authMW.Middleware)

		r.Post("/api/v1/auth/logout", h.Logout)
		r.Post("/api/v1/urls", h.CreateURL)
		r.Get("/api/v1/urls", h.ListURL)
		r.Delete("/api/v1/urls/{code}", h.DeleteURL)
		r.Get("/api/v1/analytics/{code}", h.GetAnalytics)
	})

	r.Get("/{code}", h.Redirect)

	return r
}
