package post

import (
	dom "GoNews/internal/domain/post"
	"context"
	"fmt"
)

var _ FindLatestContract = (*FindLatestUseCase)(nil)

// FindLatestUseCase представляет структуру, реализующую бизнес-логику для поиска последних n новостей.
type FindLatestUseCase struct {
	repo dom.Repository
}

// NewFindLatestUseCase создает новый экземпляр usecase для поиска последних n новостей.
func NewFindLatestUseCase(repo dom.Repository) *FindLatestUseCase {
	return &FindLatestUseCase{repo: repo}
}

// Execute выполняет бизнес-логику поиска последних n новостей.
func (uc *FindLatestUseCase) Execute(ctx context.Context, in FindLatestInputDTO) ([]PostDTO, error) {
	in.Validate()

	posts, err := uc.repo.FindLatest(ctx, in.Limit)
	if err != nil {
		return []PostDTO{}, fmt.Errorf("FindLatestUseCase.Execute: %w", err)
	}

	return MapPostsToDTO(posts), nil
}
