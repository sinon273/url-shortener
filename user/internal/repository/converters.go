package user_postgres_repository

import (
	session_domain "github.com/sinon273/url-shortener/user/internal/domain/session"
	user_domain "github.com/sinon273/url-shortener/user/internal/domain/user"
)

func userModelToDomain(model UserModel) user_domain.User {
	return user_domain.NewUser(
		model.ID,
		model.PasswordHash,
		model.Email,
	)
}

func sessionModelToDomain(model SessionModel) session_domain.Session {
	return session_domain.Session{
		ID:        model.ID,
		UserID:    model.UserID,
		ExpiresAt: model.ExpiresAt,
		IsRevoked: model.IsRevoked,
	}
}
