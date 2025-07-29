package post

import (
	dom "GoNews/internal/domain/post"
	"context"
	"fmt"
)

// Parser — интерфейс RSS-парсера.
type Parser interface {
	Parse(url string) ([]ParsedRSSDTO, error)
}

// ParseAndStoreUseCase интерфейс для парснига и сохранения RSS.
type ParseAndStoreUseCase interface {
	Execute(ctx context.Context, in ParseAndStoreInputDTO) error
}

type parseAndStoreUseCase struct {
	repo   dom.Repository
	parser Parser
}

// NewParseAndStoreUseCase создает новый экземпляр usecase для парсинга и сохранения RSS.
func NewParseAndStoreUseCase(repo dom.Repository, parser Parser) ParseAndStoreUseCase {
	return &parseAndStoreUseCase{repo: repo, parser: parser}
}

// Execute выполняет парсинг RSS ленты по указанному URL и сохраняет полученные посты в репозиторий.
func (uc *parseAndStoreUseCase) Execute(ctx context.Context, in ParseAndStoreInputDTO) error {
	if err := in.Validate(); err != nil {
		return fmt.Errorf("ParseAndStoreUseCase.Validate: %w", err)
	}

	items, err := uc.parser.Parse(in.URL)
	if err != nil {
		return fmt.Errorf("ParseAndStoreUseCase.Parse: %w", err)
	}

	for _, item := range items {
		post, err := dom.NewPost(item.Title, item.Content, item.Link, item.PubTime)
		if err != nil {
			return fmt.Errorf("ParseAndStoreUseCase.NewPost: %w", err)
		}

		existing, err := uc.repo.FindByID(ctx, post.ID())
		if err == nil && existing != nil {
			continue
		}

		if err = uc.repo.Store(ctx, post); err != nil {
			return fmt.Errorf("ParseAndStoreUseCase.Store: %w", err)
		}
	}

	return nil
}
