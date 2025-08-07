package post

import (
	"context"
	"fmt"
	"strings"

	dom "github.com/ee-crocush/go-news/go-news/internal/domain/post"
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
		Content: replaceUnnecessary(post.Content().Value()),
		Link:    post.Link().Value(),
		PubTime: post.PubTime().String(),
	}, nil
}

const UnnecessaryWords = "Читать далее"

func replaceUnnecessary(original string) string {
	clean := strings.Replace(original, UnnecessaryWords, "", 1)

	return strings.TrimSpace(clean)
}
