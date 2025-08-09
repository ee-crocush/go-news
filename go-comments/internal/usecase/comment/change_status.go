package comment

import (
	"context"
	"fmt"
	dom "github.com/ee-crocush/go-news/go-comments/internal/domain/comment"
)

var _ ChangeStatusContract = (*ChangeUseCase)(nil)

// ChangeUseCase представляет структуру, реализующую бизнес-логику для изменения статуса комментария.
type ChangeUseCase struct {
	repo dom.Repository
}

// NewChangeUseCase создает новый экземпляр usecase для изменения статуса комментария.
func NewChangeUseCase(repo dom.Repository) *ChangeUseCase {
	return &ChangeUseCase{repo: repo}
}

// Execute выполняет бизнес-логику изменения статуса комментария.
func (uc *ChangeUseCase) Execute(ctx context.Context, in StatusDTO) error {
	commentID, err := dom.NewID(in.ID)
	if err != nil {
		return fmt.Errorf("ChangeUseCase.NewID: %w", err)
	}

	status, err := dom.NewStatus(in.Status)
	if err != nil {
		return fmt.Errorf("ChangeUseCase.NewStatus: %w", err)
	}

	comment, err := uc.repo.FindByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("ChangeUseCase.FindByID: %w", err)
	}

	comment.SetStatus(status)

	if comment.IsApproved() {
		now := dom.NewTime()
		err = uc.repo.UpdateStatus(ctx, comment.ID(), comment.Status(), &now)
	} else {
		err = uc.repo.UpdateStatus(ctx, comment.ID(), comment.Status(), nil)
	}

	return nil
}
