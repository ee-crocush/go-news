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

// FromJSON сериализует событие на модерацию комментария.
func (r *CommentCreatedEvent) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

// CommentModerationResult - результат модерации комментария.
type CommentModerationResult struct {
	CommentID   int64     `json:"comment_id"`
	Status      string    `json:"status"` // "approved" или "rejected"
	ProcessedAt time.Time `json:"processed_at"`
}

// ToJSON сериализует результат модерации комментария.
func (r *CommentModerationResult) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
