package comment

import "context"

// Creator определяет контракт создания комментария.
type Creator interface {
	// Create создает комментарий.
	Create(ctx context.Context, comment *Comment) error
}

// Finder определяет контракт получения комментариев.
type Finder interface {
	// FindByID получает комментарий по ID.
	FindByID(ctx context.Context, id ID) (*Comment, error)
	// FindAll получает все комментарии.
	FindAll(ctx context.Context) ([]*Comment, error)
}
