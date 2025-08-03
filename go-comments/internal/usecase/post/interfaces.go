package post

import "context"

// FindByIDUContract интерфейс для поиска новости по ID.
type FindByIDUContract interface {
	Execute(ctx context.Context, in FindByIDInputDTO) (PostDTO, error)
}

// FindAllContract интерфейс для поиска всех новостей.
type FindAllContract interface {
	Execute(ctx context.Context) ([]PostDTO, error)
}

// FindLastContract интерфейс для поиска последней новости.
type FindLastContract interface {
	Execute(ctx context.Context) (PostDTO, error)
}

// FindLatestContract интерфейс для поиска последних n новостей.
type FindLatestContract interface {
	Execute(ctx context.Context, in FindLatestInputDTO) ([]PostDTO, error)
}
