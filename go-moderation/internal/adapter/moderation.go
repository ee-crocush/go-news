package adapter

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/go-moderation/internal/infrastructure/events"
	"github.com/ee-crocush/go-news/go-moderation/internal/service"
	"github.com/segmentio/kafka-go"
)

// ModerationAdapter инкапсулирует сервис модерации.
type ModerationAdapter struct {
	service *service.ModerationService
}

// NewModerationAdapter создает новый экземпляр ModerationAdapter.
func NewModerationAdapter(s *service.ModerationService) *ModerationAdapter {
	return &ModerationAdapter{service: s}
}

// Execute обрабатывает входящее сообщение Kafka для модерации комментария.
func (m *ModerationAdapter) Execute(ctx context.Context, msg kafka.Message) error {
	var event events.CommentCreatedEvent
	if err := event.FromJSON(msg.Value); err != nil {
		return fmt.Errorf("failed to unmarshal comment created event: %w", err)
	}

	return m.service.Moderate(ctx, event)
}
