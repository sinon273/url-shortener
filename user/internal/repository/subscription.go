package user_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sinon273/url-shortener/user/internal/domain"
)

func (r *Repository) GetSubscriptionLimit(ctx context.Context, userID string) (int, error) {
	query := `
	SELECT links_limit FROM subscriptions WHERE user_id = $1
	`

	var linksLimit int
	row := r.pool.QueryRow(ctx, query, userID)
	err := row.Scan(&linksLimit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrNotFound
		}
		return 0, fmt.Errorf("links limit: %w", err)
	}

	return linksLimit, nil
}
