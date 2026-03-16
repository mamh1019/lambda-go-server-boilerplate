package kafka

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	brokerAddr := os.Getenv("KAFKA_BROKER_ADDR")
	if brokerAddr == "" {
		brokerAddr = "localhost:9092"
	}

	topic := os.Getenv("KAFKA_USER_CREATED_TOPIC")
	if topic == "" {
		topic = "user-created"
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(brokerAddr),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}

	return &Producer{writer: w}
}

func (p *Producer) PublishUserCreated(ctx context.Context, userID int64) error {
	if p == nil || p.writer == nil {
		return nil
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", userID)),
		Value: []byte(fmt.Sprintf(`{"user_id":%d}`, userID)),
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(ctx, msg)
}
