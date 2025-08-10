package kafka

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/events"
	"github.com/ee-crocush/go-news/pkg/logger"
	"github.com/segmentio/kafka-go"
)

// ModerationResultProcessor интерфейс для обработки результатов модерации.
type ModerationResultProcessor interface {
	Execute(ctx context.Context, result events.CommentModerationResult) error
}

type Consumer struct {
	reader    *kafka.Reader
	processor ModerationResultProcessor
}

func NewConsumer(brokers []string, topic, groupID string, processor ModerationResultProcessor) *Consumer {
	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			CommitInterval: 0,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
		},
	)

	return &Consumer{reader: reader, processor: processor}
}

// Start запускает consumer для обработки сообщений.
func (c *Consumer) Start(ctx context.Context) error {
	log := logger.GetLogger()
	log.Info().Msg("Starting Kafka consumer for moderation results...")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Stopping Kafka consumer...")
			return c.reader.Close()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				return fmt.Errorf("failed to fetch message: %w", err)
			}

			var result events.CommentModerationResult
			if err = result.FromJSON(msg.Value); err != nil {
				return fmt.Errorf("failed to parse moderation result: %w", err)
			}

			if err = c.processor.Execute(ctx, result); err != nil {
				log.
					Err(err).
					Str("comment_id", fmt.Sprintf("%d", result.CommentID)).
					Str("comment_status", result.Status).
					Msg("Failed to to update comment status")
				continue
			}

			if err = c.reader.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("failed to commit message: %w", err)
			}

			log.Info().Msgf("Successfully update status=%s for comment ID=%d", result.Status, result.CommentID)
		}
	}
}
