package url_producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
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
	return nil

}
