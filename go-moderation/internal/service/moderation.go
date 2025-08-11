// Package service представляет работу сервиса модерации
package service

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/go-moderation/internal/infrastructure/events"
	"strings"
	"time"
)

// ModerationService представляет сервис модерации.
type ModerationService struct {
	publisher Publisher
}

// Publisher интерфейс, который передаем в кафку для выполнения событий.
type Publisher interface {
	Publish(ctx context.Context, key string, value []byte) error
}

// NewService создает новый экземпляр Service.
func NewService(pub Publisher) *ModerationService {
	return &ModerationService{publisher: pub}
}

// bannedPhrases представляет "запрещенные" слова
var bannedPhrases = []string{
	"zxcvbn",
	"qwerty",
	"asdfgh",
	"йцукен",
	"фывапр",
	"ячсмит",
}

const (
	Approved = "approved"
	Rejected = "rejected"
)

// Status представляет удобную обертку статуса.
// Можно было бы реализовать в виде флага, но так более наглядно.
type Status struct {
	Value string
}

// Moderate выполняет модерацию контента.
func (s *ModerationService) Moderate(ctx context.Context, e events.CommentCreatedEvent) error {
	status := Status{Value: Approved}
	contentLower := strings.ToLower(e.Content)

	for _, phrase := range bannedPhrases {
		if strings.Contains(contentLower, phrase) {
			status.Value = Rejected
			break
		}
	}

	result := events.CommentModerationResult{
		CommentID:   e.CommentID,
		Status:      status.Value,
		ProcessedAt: time.Now(),
	}

	data, err := result.ToJSON()
	if err != nil {
		return fmt.Errorf("Moderate.ToJSON: %w", err)
	}

	return s.publisher.Publish(ctx, fmt.Sprintf("%d", e.CommentID), data)
}
