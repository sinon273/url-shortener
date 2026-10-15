package analytics_consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
)

func (c *Consumer) Run(ctx context.Context) error {
	msgCh := make(chan kafka.Message)
	errCh := make(chan error, 1)

	go func() {
		for {
			msg, err := c.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				errCh <- err
				return
			}
			msgCh <- msg
		}
	}()
	timer := time.NewTimer(c.flushInterval)
	defer timer.Stop()
	for {
		select {
		case msg := <-msgCh:
			var msgJson analytics_domain.ClickEvent
			err := json.Unmarshal(msg.Value, &msgJson)
			if err != nil {
				return fmt.Errorf("unmarshal analytics click event: %w", err)
			}
			c.events = append(c.events, msgJson)
			c.msgs = append(c.msgs, msg)

			if len(c.events) >= c.batchSize {
				if err = c.flush(ctx); err != nil {
					return fmt.Errorf("flush analytics click events: %w", err)
				}
				timer.Reset(c.flushInterval)
			}
		case <-timer.C:
			err := c.flush(ctx)
			if err != nil {
				return fmt.Errorf("flush: %w", err)
			}
			timer.Reset(c.flushInterval)
		case err := <-errCh:
			if err != nil {
				return fmt.Errorf("fetch: %w", err)
			}
		case <-ctx.Done():
			if len(c.events) > 0 {
				flushCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				err := c.flush(flushCtx)
				if err != nil {
					return fmt.Errorf("flush: %w", err)
				}
			}
			return nil
		}
	}
}
