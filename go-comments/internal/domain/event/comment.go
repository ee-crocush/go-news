// Package event содержит события на создание и изменение комментариев.
package event

import (
	"encoding/json"
	"time"
)

// CommentCreatedEvent - событие создания комментария для модерации.
type CommentCreatedEvent struct {
	CommentID int64     `json:"comment_id"`
	NewsID    int32     `json:"news_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ToJSON конвертирует событие в JSON.
func (e *CommentCreatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// CommentModerationResult - результат модерации комментария.
type CommentModerationResult struct {
	CommentID   int64     `json:"comment_id"`
	Status      string    `json:"status"`           // "approved" или "rejected"
	Reason      string    `json:"reason,omitempty"` // Причина отклонения
	ProcessedAt time.Time `json:"processed_at"`
}

// FromJSON создает результат модерации из JSON.
func (r *CommentModerationResult) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

// Константы для топиков Kafka.
const (
	TopicCommentCreated          = "comment.created"
	TopicCommentModerationResult = "comment.moderation.result"
)
