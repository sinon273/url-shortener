package analytics_consumer

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
)

type Consumer struct {
	*kafka.Reader
	repo          Repository
	batchSize     int
	flushInterval time.Duration
	events        []analytics_domain.ClickEvent
	msgs          []kafka.Message
}

func NewConsumer(reader *kafka.Reader, repo Repository, batchSize int, flushInterval time.Duration) *Consumer {
	return &Consumer{
		Reader:        reader,
		repo:          repo,
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}
}

type Repository interface {
	InsertBatch(ctx context.Context, events []analytics_domain.ClickEvent) error
}
