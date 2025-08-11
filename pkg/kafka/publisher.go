package kafka

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/pkg/logger"
	"time"

	"github.com/segmentio/kafka-go"
)

// Publisher представляет Kafka publisher.
type Publisher struct {
	writer *kafka.Writer
}

// NewPublisher создает новый Kafka Publisher.
func NewPublisher(brokers []string, topic string) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

// Publish отправляет событие.
func (p *Publisher) Publish(ctx context.Context, key string, value []byte) error {
	log := logger.GetLogger()
	log.Info().Msgf("Starting publishing message key: %s", key)

	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	log.Info().Msgf("Message published successfully key: %s", key)

	return nil
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
