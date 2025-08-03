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
	// FindAllByNewsID получает все комментарии для конкрентой новости.
	FindAllByNewsID(ctx context.Context, newsID NewsID) ([]*Comment, error)
}
