package analytics_repository

import (
	"context"
	"fmt"
)

func (r *Repository) EnsureTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS click_events (
	short_code String,
	user_id    String,
	clicked_at DateTime,
	ip	       String,
	user_agent String,
	referer    String
) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(clicked_at)
	ORDER BY (short_code, clicked_at)
	`

	err := r.conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("create click_events table: %w", err)
	}
	return nil
}
