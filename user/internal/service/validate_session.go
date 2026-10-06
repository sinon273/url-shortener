package user_service

import (
	"context"

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
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return "", "", domain.ErrUnauthenticated
	}
	if !session.IsAlive() {
		return "", "", domain.ErrUnauthenticated
	}

	_ = s.cache.SetSession(ctx, sessionID, session.UserID, s.auth.AccessTTL())

	return session.UserID, sessionID, nil
}
