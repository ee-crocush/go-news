// Package events содержит события на создание и изменение комментариев.
package events

import (
	"encoding/json"
	"time"
)

// CommentCreatedEvent - событие создания комментария для модерации.
type CommentCreatedEvent struct {
	CommentID int64     `json:"comment_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// NewCommentCreatedEvent создает экземпляр CommentCreatedEvent.
func NewCommentCreatedEvent(commentID int64, content string) *CommentCreatedEvent {
	return &CommentCreatedEvent{
		CommentID: commentID,
		Content:   content,
		CreatedAt: time.Now(),
	}
}

// ToJSON конвертирует событие в JSON.
func (e *CommentCreatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// CommentModerationResult - результат модерации комментария.
type CommentModerationResult struct {
	CommentID   int64     `json:"comment_id"`
	Status      string    `json:"status"` // "approved" или "rejected"
	ProcessedAt time.Time `json:"processed_at"`
}

// FromJSON создает результат модерации из JSON.
func (r *CommentModerationResult) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}
