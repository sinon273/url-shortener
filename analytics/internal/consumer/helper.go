package analytics_consumer

import (
	"context"
	"fmt"
	"log/slog"

	analytics_metrics "github.com/sinon273/url-shortener/analytics/internal/metrics"
)

func (c *Consumer) flush(ctx context.Context) error {
	err := c.repo.InsertBatch(ctx, c.events)
	if err != nil {
		slog.ErrorContext(ctx, "failed to insert batch", "err", err)
		return fmt.Errorf("insert batch: %w", err)
	}
	analytics_metrics.KafkaMessagesConsumedTotal.Add(float64(len(c.events)))
	err = c.CommitMessages(ctx, c.msgs...)
	if err != nil {
		slog.ErrorContext(ctx, "failed to commit messages", "err", err)
		c.events = c.events[:0]
		c.msgs = c.msgs[:0]
		return fmt.Errorf("commit messages: %w", err)
	}
	c.events = c.events[:0]
	c.msgs = c.msgs[:0]

	return nil
}
