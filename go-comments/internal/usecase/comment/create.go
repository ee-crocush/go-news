package comment

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/pkg/logger"
	"time"

	dom "github.com/ee-crocush/go-news/go-comments/internal/domain/comment"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/events"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/kafka"
)

var _ CreateContract = (*CreateUseCase)(nil)

// CreateUseCase представляет структуру, реализующую бизнес-логику для создания комментария.
type CreateUseCase struct {
	repo      dom.Repository
	publisher kafka.EventPublisher
}

// NewCreateUseCase создает новый экземпляр usecase для создания комментария.
func NewCreateUseCase(repo dom.Repository, publisher kafka.EventPublisher) *CreateUseCase {
	return &CreateUseCase{repo: repo, publisher: publisher}
}

// Execute выполняет бизнес-логику создания комментария.
func (uc *CreateUseCase) Execute(ctx context.Context, in CommentDTO) error {
	comment, err := dom.NewComment(in.NewsID, in.Username, in.Content)
	if err != nil {
		return fmt.Errorf("CreateUseCase.NewComment: %w", err)
	}

	if in.ParentID != nil {
		parentID, err := dom.NewParentID(*in.ParentID)
		if err != nil {
			return fmt.Errorf("CreateUseCase.NewParentID: %w", err)
		}

		existingID, _ := dom.NewID(*in.ParentID)
		if _, err = uc.repo.FindByID(ctx, existingID); err != nil {
			return fmt.Errorf("CreateUseCase.FindByID: %w", err)
		}

		comment.SetParentID(parentID)
	}

	commentID, err := uc.repo.Create(ctx, comment)
	if err != nil {
		return fmt.Errorf("CreateUseCase.Create: %w", err)
	}
	comment.SetID(commentID)

	event := &events.CommentCreatedEvent{
		CommentID: comment.ID().Value(),
		Content:   comment.Content().Value(),
		CreatedAt: time.Now(),
	}
	// Логгируем тут, чтобы не пропустить косяк
	if err = uc.publisher.PublishCommentCreated(ctx, event); err != nil {
		log := logger.GetLogger()
		log.
			Err(err).
			Str("comment_id", fmt.Sprintf("%d", commentID.Value())).
			Msg("Failed to publish comment created event")
	}

	return nil
}
