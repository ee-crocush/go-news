// Package comment выполняет бизнес-логику по комментариям.
package comment

import (
	"context"
	"github.com/segmentio/kafka-go"
)

// CreateContract интерфейс для создания комментария.
type CreateContract interface {
	Execute(ctx context.Context, in CommentDTO) error
}

// FindAllByNewsIDContract интерфейс для поиска всех комментариев для конкретной новости.
type FindAllByNewsIDContract interface {
	Execute(ctx context.Context, in AllByNewsIDDTO) ([]CommentDTO, error)
}

// ChangeStatusContract интерфейс для публикации/отклонения комментария.
type ChangeStatusContract interface {
	Execute(ctx context.Context, msg kafka.Message) error
}
