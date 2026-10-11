package gateway_middleware

import (
	"context"
	"net/http"
	"strings"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

type Auth struct {
	userClient userv1.UserServiceClient
}

func NewAuth(userClient userv1.UserServiceClient) *Auth {
	return &Auth{userClient: userClient}
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		req := &userv1.ValidateSessionRequest{
			AccessToken: token,
		}
		resp, err := a.userClient.ValidateSession(r.Context(), req)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		userID, sessionID := resp.UserId, resp.SessionId
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, sessionIDKey, sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
