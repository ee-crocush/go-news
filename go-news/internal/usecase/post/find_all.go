package post

import (
	"context"
	"fmt"
	dom "github.com/ee-crocush/go-news/go-news/internal/domain/post"
)

var _ FindAllContract = (*FindAllUseCase)(nil)

// FindAllUseCase представляет структуру, реализующую бизнес-логику для поиска всех новостей.
type FindAllUseCase struct {
	repo dom.Repository
}

// NewFindAllUseCase создает новый экземпляр adapter для поиска всех новостей.
func NewFindAllUseCase(repo dom.Repository) *FindAllUseCase {
	return &FindAllUseCase{repo: repo}
}

// Execute выполняет бизнес-логику поиска всех новостей.
func (uc *FindAllUseCase) Execute(ctx context.Context, in FindAllInputDTO) ([]PostDTO, int32, error) {
	offset := (in.Page - 1) * in.Limit

	posts, total, err := uc.repo.FindAll(ctx, in.Search, in.Limit, offset)
	if err != nil {
		return []PostDTO{}, 0, fmt.Errorf("FindAllUseCase.Execute: %w", err)
	}

	return MapPostsToDTO(posts), total, nil
}
