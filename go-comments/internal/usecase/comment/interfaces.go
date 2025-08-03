package comment

import "context"

// CreateContract интерфейс для создания комментария.
type CreateContract interface {
	Execute(ctx context.Context, in CommentDTO) error
}

// FindAllByNewsIDContract интерфейс для поиска всех комментариев для конкретной новости.
type FindAllByNewsIDContract interface {
	Execute(ctx context.Context, in AllByNewsIDDTO) ([]CommentDTO, error)
}
