package user_service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sinon273/url-shortener/user/internal/domain"
)

func (s *Service) Login(ctx context.Context, email string, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", domain.ErrUnauthenticated
	}
	err = s.auth.ComparePassword(user.PasswordHash, password)
	if err != nil {
		return "", "", domain.ErrUnauthenticated
	}

	expiresAt := time.Now().Add(s.auth.RefreshTTL())
	sessionID, err := s.repo.CreateSession(ctx, user.ID, expiresAt)
	if err != nil {
		return "", "", fmt.Errorf("create session: %w", err)
	}
	access, err := s.auth.SignAccess(user.ID, sessionID)
	if err != nil {
		return "", "", fmt.Errorf("sign access: %w", err)
	}
	refresh, err := s.auth.SignRefresh(user.ID, sessionID)
	if err != nil {
		return "", "", fmt.Errorf("sign refresh: %w", err)
	}

	err = s.cache.SetSession(ctx, sessionID, user.ID, s.auth.AccessTTL())
	if err != nil {
		slog.ErrorContext(ctx, "failed to set session cache", "err", err, "session_id", sessionID)
	}
	return access, refresh, nil
}
