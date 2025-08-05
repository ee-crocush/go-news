package post

import (
	"context"
	"fmt"
	dom "github.com/ee-crocush/go-news/go-news/internal/domain/post"
)

var _ FindByTitleSubstringContract = (*FindByTitleSubstring)(nil)

// FindByTitleSubstring представляет структуру, реализующую бизнес-логику для поиска всех новостей.
type FindByTitleSubstring struct {
	repo dom.Repository
}

// NewFindByTitleSubstring создает новый экземпляр usecase для поиска всех новостей.
func NewFindByTitleSubstring(repo dom.Repository) *FindByTitleSubstring {
	return &FindByTitleSubstring{repo: repo}
}

// Execute выполняет бизнес-логику поиска всех новостей.
func (uc *FindByTitleSubstring) Execute(ctx context.Context, in FindByTitleSubstringInputDTO) ([]PostDTO, error) {
	posts, err := uc.repo.FindByTitleSubstring(ctx, in.Substring)
	if err != nil {
		return []PostDTO{}, fmt.Errorf("FindByTitleSubstring.Execute: %w", err)
	}

	return MapPostsToDTO(posts), nil
}
