package post

import (
	"context"
	"fmt"
	dom "github.com/ee-crocush/go-news/go-news/internal/domain/post"
)

var _ FindLastContract = (*FindLastUseCase)(nil)

// FindLastUseCase представляет структуру, реализующую бизнес-логику для поиска последней новости.
type FindLastUseCase struct {
	repo dom.Repository
}

// NewFindLastUseCase создает новый экземпляр adapter для поиска последней новости.
func NewFindLastUseCase(repo dom.Repository) *FindLastUseCase {
	return &FindLastUseCase{repo: repo}
}

func (uc *FindLastUseCase) Execute(ctx context.Context) (PostDTO, error) {
	post, err := uc.repo.FindLast(ctx)
	if err != nil {
		return PostDTO{}, fmt.Errorf("FindLastUseCase.Execute: %w", err)
	}

	return PostDTO{
		ID:      post.ID().Value(),
		Title:   post.Title().Value(),
		Content: post.Content().Value(),
		Link:    post.Link().Value(),
		PubTime: post.PubTime().String(),
	}, nil
}
