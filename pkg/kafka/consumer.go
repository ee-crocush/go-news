package kafka

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/pkg/logger"

	"github.com/segmentio/kafka-go"
)

// ConsumerProcessor интерфейс для обработки результатов.
type ConsumerProcessor interface {
	Execute(ctx context.Context, msg kafka.Message) error
}

// Consumer представляет собой Kafka consumer.
type Consumer struct {
	reader  *kafka.Reader
	handler ConsumerProcessor
}

// NewConsumer создает новый экземпляр Consumer.
func NewConsumer(brokers []string, topic, groupID string, handler ConsumerProcessor) *Consumer {
	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			CommitInterval: 0,
		},
	)
	return &Consumer{reader: reader, handler: handler}
}

// Start запускает consumer для обработки сообщений.
func (c *Consumer) Start(ctx context.Context) error {
	log := logger.GetLogger()
	log.Info().Msg("Starting Kafka consumer...")

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

			if err = c.handler.Execute(ctx, msg); err != nil {
				log.
					Err(err).
					RawJSON(string(msg.Key), msg.Value).
					Msg("Failed to to update comment status")
				continue
			}

			if err = c.reader.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("failed to commit message: %w", err)
			}

			log.Info().RawJSON(string(msg.Key), msg.Value).Msg("Kafka message processed successfully")
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
