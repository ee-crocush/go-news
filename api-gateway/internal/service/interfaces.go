package service

import (
	"context"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/dto"
)

// NewsContract интерфейс для работы с новостями
type NewsContract interface {
	GetAllNews(ctx context.Context) ([]dto.Post, error)
	GetLastNews(ctx context.Context) (*dto.Post, error)
	GetLatestNews(ctx context.Context, limit int) ([]dto.Post, error)
	GetNewsByID(ctx context.Context, id int32) (*dto.PostWithComments, error)
}

// CommentContract интерфейс для работы с комментариями
type CommentContract interface {
	GetCommentsByNewsID(ctx context.Context, newsID int32) ([]dto.Comment, error)
	CreateComment(ctx context.Context, req dto.CreateCommentRequest) error
}
