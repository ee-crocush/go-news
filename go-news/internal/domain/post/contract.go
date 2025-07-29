package post

import "context"

// PostStore определяет контракт сохранения новости.
type PostStore interface {
	// Store сохраняет новость.
	Store(ctx context.Context, post *Post) error
}

// PostFinder определяет контракт получения новостей.
type PostFinder interface {
	// FindByID получает новость по ID.
	FindByID(ctx context.Context, postID PostID) (*Post, error)
	// FindLast получает последнюю новость.
	FindLast(ctx context.Context) (*Post, error)
	// FindLatest получает последние n новостей.
	FindLatest(ctx context.Context, limit int) ([]*Post, error)
	// FindAll получает все новости.
	FindAll(ctx context.Context) ([]*Post, error)
}
