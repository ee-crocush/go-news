package post

import (
	"context"
	"fmt"

	dom "GoNews/internal/domain/post"
)

var _ FindByIDUContract = (*FindByIDUseCase)(nil)

// FindByIDUseCase представляет структуру, реализующую бизнес-логику для поиска новости по ID.
type FindByIDUseCase struct {
	repo dom.Repository
}

// NewFindByIDUseCase создает новый экземпляр usecase для поиска новости по ID.
func NewFindByIDUseCase(repo dom.Repository) *FindByIDUseCase {
	return &FindByIDUseCase{repo: repo}
}

// Execute выполняет бизнес-логику поиска новости по ID.
func (uc *FindByIDUseCase) Execute(ctx context.Context, in FindByIDInputDTO) (PostDTO, error) {
	postID, err := dom.NewPostID(in.ID)
	if err != nil {
		return PostDTO{}, fmt.Errorf("FindByIDUseCase.NewPostID: %w", err)
	}

	post, err := uc.repo.FindByID(ctx, postID)
	if err != nil {
		return PostDTO{}, fmt.Errorf("findByIDUseCase.FindByID: %w", err)
	}

	return PostDTO{
		ID:      post.ID().Value(),
		Title:   post.Title().Value(),
		Content: post.Content().Value(),
		Link:    post.Link().Value(),
		PubTime: post.PubTime().String(),
	}, nil
}
