// Package kafka представляет пакет для публикации и чтении сообщений из топиков.
package kafka

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/events"
	"github.com/ee-crocush/go-news/pkg/logger"
	"github.com/segmentio/kafka-go"
	"time"
)

// EventPublisher интерфейс для публикации событий.
type EventPublisher interface {
	PublishCommentCreated(ctx context.Context, event *events.CommentCreatedEvent) error
}

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

// PublishCommentCreated отправляет событие о создании комментария.
func (p *Publisher) PublishCommentCreated(ctx context.Context, event *events.CommentCreatedEvent) error {
	log := logger.GetLogger()
	log.Info().Msgf("Starting Kafka publisher for moderate comment with ID: %d", event.CommentID)

	data, err := event.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize events: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", event.CommentID)),
		Value: data,
		Time:  time.Now(),
	}
	if err = p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	log.Info().Msgf("Successfully published comment created event for comment ID: %d", event.CommentID)

	return nil
}

// Close закрывает writer.
func (p *Publisher) Close() error {
	return p.writer.Close()
}
