package url_producer

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type ClickEvent struct {
	ShortCode string    `json:"short_code"`
	UserID    string    `json:"user_id"`
	ClickedAt time.Time `json:"clicked_at"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(writer *kafka.Writer) *Producer {
	return &Producer{writer: writer}
}
