package analytics_repository

import (
	"context"
	"fmt"

	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
)

func (r *Repository) InsertBatch(ctx context.Context, events []analytics_domain.ClickEvent) error {
	if len(events) == 0 {
		return nil
	}
	query := `
	INSERT INTO click_events (short_code, user_id, clicked_at, ip, user_agent, referer)
	`
	batch, err := r.conn.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}
	for _, event := range events {
		err = batch.Append(event.ShortCode, event.UserID, event.ClickedAt, event.IP, event.UserAgent, event.Referer)
		if err != nil {
			return fmt.Errorf("append event: %w", err)
		}
	}
	err = batch.Send()
	if err != nil {
		return fmt.Errorf("send batch: %w", err)
	}
	return nil
}
