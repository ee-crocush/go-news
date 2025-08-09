// Package domain содержит выходные данные запросов
package dto

// CreateCommentRequest представляет тело запроса для создания комментария.
type CreateCommentRequest struct {
	NewsID   int32  `json:"news_id" example:"1"`
	ParentID *int64 `json:"parent_id" example:"1"`
	Username string `json:"username" example:"Example_username"`
	Content  string `json:"content" example:"Example content"`
}

// Comment описывает структуру комментария.
type Comment struct {
	ID       int64     `json:"id" example:"1"`
	NewsID   int32     `json:"news_id" example:"1"`
	ParentID *int64    `json:"parent_id" example:"1"`
	Username string    `json:"username" example:"Example_username"`
	Content  string    `json:"content" example:"Example content"`
	PubTime  string    `json:"pub_time" example:"2025-06-26 10:00:43"`
	Children []Comment `json:"children" example:"[]"`
}

// CommentResponse описывает структуру ответа на /comments.
type CommentResponse struct {
	Status  string `json:"status" example:"OK"`
	Message string `json:"message" example:"Example message success"`
}
