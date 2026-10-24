package user_service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/sinon273/url-shortener/user/internal/domain"
)

func (s *Service) ValidateSession(ctx context.Context, accessToken string) (userID, sessionID string, err error) {
	userID, sessionID, err = s.auth.Parse(accessToken)
	if err != nil {
		return "", "", domain.ErrUnauthenticated
	}

	userID, err = s.cache.GetSession(ctx, sessionID)
	if err == nil {
		return userID, sessionID, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		slog.ErrorContext(ctx, "failed to get session cache", "err", err, "session_id", sessionID)
	}

	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return "", "", domain.ErrUnauthenticated
	}
	if !session.IsAlive() {
		return "", "", domain.ErrUnauthenticated
	}

	err = s.cache.SetSession(ctx, sessionID, session.UserID, s.auth.AccessTTL())
	if err != nil {
		slog.ErrorContext(ctx, "failed to set session cache", "err", err, "session_id", sessionID)
	}
	return session.UserID, sessionID, nil
}
