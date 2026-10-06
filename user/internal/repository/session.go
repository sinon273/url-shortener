package user_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sinon273/url-shortener/user/internal/domain"
	session_domain "github.com/sinon273/url-shortener/user/internal/domain/session"
)

func (r *Repository) CreateSession(ctx context.Context, userID string, expiresAt time.Time) (string, error) {
	query := `
	INSERT INTO sessions (user_id, expires_at) VALUES($1,$2) RETURNING id
	`

	var sessionID string
	row := r.pool.QueryRow(ctx, query, userID, expiresAt)
	err := row.Scan(&sessionID)
	if err != nil {
		return "", fmt.Errorf("insert session: %w", err)
	}

	return sessionID, nil
}

func (r *Repository) GetSession(ctx context.Context, sessionID string) (session_domain.Session, error) {
	query := `
	SELECT id, user_id, expires_at, is_revoked FROM sessions WHERE id = $1
	`
	var sessionModel SessionModel
	row := r.pool.QueryRow(ctx, query, sessionID)
	err := row.Scan(
		&sessionModel.ID,
		&sessionModel.UserID,
		&sessionModel.ExpiresAt,
		&sessionModel.IsRevoked,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return session_domain.Session{}, domain.ErrNotFound
		}
		return session_domain.Session{}, fmt.Errorf("get session: %w", err)
	}

	return sessionModelToDomain(sessionModel), nil
}

func (r *Repository) RevokeSession(ctx context.Context, sessionID string) error {
	query := `
	UPDATE sessions SET is_revoked = TRUE, updated_at = NOW() WHERE id = $1
	`

	tag, err := r.pool.Exec(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("exec update session: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
