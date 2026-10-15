package analytics_consumer

import (
	"context"
	"fmt"
)

func (c *Consumer) flush(ctx context.Context) error {
	err := c.repo.InsertBatch(ctx, c.events)
	if err != nil {
		return fmt.Errorf("insert batch: %w", err)
	}
	err = c.CommitMessages(ctx, c.msgs...)
	if err != nil {
		c.events = c.events[:0]
		c.msgs = c.msgs[:0]
		return fmt.Errorf("commit messages: %w", err)
	}
	c.events = c.events[:0]
	c.msgs = c.msgs[:0]

	return nil
}
