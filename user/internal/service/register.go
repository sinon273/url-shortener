package user_service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/sinon273/url-shortener/user/internal/domain"
)

func (s *Service) Register(
	ctx context.Context,
	email,
	password string,
) (string, error) {

	if _, err := mail.ParseAddress(email); err != nil {
		return "", domain.ErrInvalidArgument
	}
	if len(password) < 8 {
		return "", domain.ErrInvalidArgument
	}
	passwordHash, err := s.auth.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("error hash password: %w", err)
	}
	userID, err := s.repo.CreateUserWithSubscription(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return "", domain.ErrEmailAlreadyExists
		}
		return "", fmt.Errorf("error create user: %w", err)
	}

	return userID, nil

}
