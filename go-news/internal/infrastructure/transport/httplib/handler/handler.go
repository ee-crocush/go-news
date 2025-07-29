// Package handler содержит все обработчики HTTP запросов
package handler

import (
	uc "GoNews/internal/usecase/post"
	"context"
)

// FindByIDPostExecutor интерфейс для поиска новости по ID.
type FindByIDPostExecutor interface {
	Execute(ctx context.Context, in uc.FindByIDInputDTO) (uc.PostDTO, error)
}

// FindLastPostExecutor интерфейс для поиска последней новости.
type FindLastPostExecutor interface {
	Execute(ctx context.Context) (uc.PostDTO, error)
}

// FindLatestPostExecutor интерфейс для поиска последних n новостей.
type FindLatestPostExecutor interface {
	Execute(ctx context.Context, in uc.FindLatestInputDTO) ([]uc.PostDTO, error)
}

// FindAllPostExecutor интерфейс для поиска всех новостей.
type FindAllPostExecutor interface {
	Execute(ctx context.Context) ([]uc.PostDTO, error)
}

// Handler представляет HTTP-handler для работы с новостями.
type Handler struct {
	findByIDUC   FindByIDPostExecutor
	findLastUC   FindLastPostExecutor
	findLatestUC FindLatestPostExecutor
	findAllUC    FindAllPostExecutor
}

// NewHandler создает новый экземпляр HTTP-handler.
func NewHandler(
	findByIDUC FindByIDPostExecutor,
	findLastUC FindLastPostExecutor,
	findLatestUC FindLatestPostExecutor,
	findAllUC FindAllPostExecutor,
) *Handler {
	return &Handler{
		findByIDUC:   findByIDUC,
		findLastUC:   findLastUC,
		findLatestUC: findLatestUC,
		findAllUC:    findAllUC,
	}
}
