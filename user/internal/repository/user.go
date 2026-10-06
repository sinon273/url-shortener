package user_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sinon273/url-shortener/user/internal/domain"
	user_domain "github.com/sinon273/url-shortener/user/internal/domain/user"
)

func (r *Repository) CreateUserWithSubscription(ctx context.Context, email, passwordHash string) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID string

	queryUser := `
	INSERT INTO users (email, password_hash) 
	VALUES ($1, $2) RETURNING id
	`

	querySubscription := `
	INSERT INTO subscriptions (user_id, plan, links_limit) VALUES($1,$2,$3)
	`

	row := tx.QueryRow(ctx, queryUser, email, passwordHash)
	err = row.Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", domain.ErrEmailAlreadyExists
		}
		return "", fmt.Errorf("insert user: %w", err)
	}

	_, err = tx.Exec(ctx, querySubscription, userID, "basic", 1000)
	if err != nil {
		return "", fmt.Errorf("insert subscription: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit transaction: %w", err)
	}

	return userID, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (user_domain.User, error) {
	query := `
	SELECT id, email, password_hash FROM users WHERE email = $1
	`

	row := r.pool.QueryRow(ctx, query, email)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user_domain.User{}, domain.ErrNotFound
		}
		return user_domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	return userModelToDomain(userModel), nil
}
