package url_producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	url_metrics "github.com/sinon273/url-shortener/url/internal/metrics"
)

func (p *Producer) PublishClick(ctx context.Context, event ClickEvent) error {
	byteEvent, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal click event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{Value: byteEvent})
	if err != nil {
		return fmt.Errorf("write click event: %w", err)
	}
	url_metrics.KafkaMessagesProduced.Inc()
	return nil

}
