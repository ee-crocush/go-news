package comment

import (
	"context"
	"fmt"
	dom "github.com/ee-crocush/go-news/go-comments/internal/domain/comment"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/events"
	"github.com/segmentio/kafka-go"
)

var _ ChangeStatusContract = (*ChangeStatusUseCase)(nil)

// ChangeStatusUseCase представляет структуру, реализующую бизнес-логику для изменения статуса комментария.
type ChangeStatusUseCase struct {
	repo dom.Repository
}

// NewChangeStatusUseCase создает новый экземпляр adapter для изменения статуса комментария.
func NewChangeStatusUseCase(repo dom.Repository) *ChangeStatusUseCase {
	return &ChangeStatusUseCase{repo: repo}
}

// Execute выполняет бизнес-логику изменения статуса комментария.
func (uc *ChangeStatusUseCase) Execute(ctx context.Context, msg kafka.Message) error {
	var in events.CommentModerationResult
	// Парсим результат из кафки
	if err := in.FromJSON(msg.Value); err != nil {
		return fmt.Errorf("ChangeUseCase.FromJSON: %w", err)
	}

	commentID, err := dom.NewID(in.CommentID)
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
