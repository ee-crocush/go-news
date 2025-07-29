package post

import (
	dom "GoNews/internal/domain/post"
	"context"
	"fmt"
)

var _ FindAllContract = (*FindAllUseCase)(nil)

// FindAllUseCase представляет структуру, реализующую бизнес-логику для поиска всех новостей.
type FindAllUseCase struct {
	repo dom.Repository
}

// NewFindAllUseCase создает новый экземпляр usecase для поиска всех новостей.
func NewFindAllUseCase(repo dom.Repository) *FindAllUseCase {
	return &FindAllUseCase{repo: repo}
}

// Execute выполняет бизнес-логику поиска всех новостей.
func (uc *FindAllUseCase) Execute(ctx context.Context) ([]PostDTO, error) {
	posts, err := uc.repo.FindAll(ctx)
	if err != nil {
		return []PostDTO{}, fmt.Errorf("FindAllUseCase.Execute: %w", err)
	}

	return MapPostsToDTO(posts), nil
}
