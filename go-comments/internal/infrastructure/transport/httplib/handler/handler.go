// Package handler содержит все обработчики HTTP запросов
package handler

import (
	"context"
	uc "github.com/ee-crocush/go-news/go-comments/internal/usecase/comment"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// CreateCommentExecutor интерфейс для создания комментария.
type CreateCommentExecutor interface {
	Execute(ctx context.Context, in uc.CommentDTO) error
}

// FindAllByNewsExecutor интерфейс для поиска всех комментариев для заданной новости.
type FindAllByNewsExecutor interface {
	Execute(ctx context.Context, in uc.AllByNewsIDDTO) ([]uc.CommentDTO, error)
}

// Handler представляет HTTP-handler для работы с комментариями.
type Handler struct {
	createUC        CreateCommentExecutor
	findAllByNewsUC FindAllByNewsExecutor
}

// NewHandler создает новый экземпляр HTTP-handler.
func NewHandler(createUC CreateCommentExecutor, findAllByNewsUC FindAllByNewsExecutor) *Handler {
	return &Handler{
		createUC:        createUC,
		findAllByNewsUC: findAllByNewsUC,
	}
}
