package post

import (
	dom "GoNews/internal/domain/post"
	"context"
	"errors"
	"testing"
)

// mockRepositoryForFindLast implements dom.Repository for testing
type mockRepositoryForFindLast struct {
	post *dom.Post
	err  error
}

func (m *mockRepositoryForFindLast) FindLast(ctx context.Context) (*dom.Post, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.post, nil
}

func (m *mockRepositoryForFindLast) FindByID(ctx context.Context, postID dom.PostID) (*dom.Post, error) {
	return nil, nil
}

func (m *mockRepositoryForFindLast) FindAll(ctx context.Context) ([]*dom.Post, error) {
	return nil, nil
}

func (m *mockRepositoryForFindLast) FindLatest(ctx context.Context, limit int) ([]*dom.Post, error) {
	return nil, nil
}

func (m *mockRepositoryForFindLast) Store(ctx context.Context, post *dom.Post) error {
	return nil
}

func TestNewFindLastUseCase(t *testing.T) {
	repo := &mockRepositoryForFindLast{}
	useCase := NewFindLastUseCase(repo)

	if useCase == nil {
		t.Error("expected useCase to not be nil")
	}

	if useCase.repo != repo {
		t.Error("expected repo to be set correctly")
	}
}

func TestFindLastUseCase_Execute(t *testing.T) {
	t.Run(
		"successful execution", func(t *testing.T) {
			postID, _ := dom.NewPostID(1)
			post, _ := dom.NewPost(
				"Last News Title", "Last News Content", "https://lastnews.com", dom.NewPubTime().Time().Unix(),
			)
			post.SetID(postID)

			repo := &mockRepositoryForFindLast{
				post: post,
				err:  nil,
			}

			useCase := NewFindLastUseCase(repo)
			ctx := context.Background()

			result, err := useCase.Execute(ctx)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if result.ID != 1 {
				t.Errorf("expected post ID to be 1, got %d", result.ID)
			}
			if result.Title != "Last News Title" {
				t.Errorf("expected post title to be 'Last News Title', got %s", result.Title)
			}
			if result.Content != "Last News Content" {
				t.Errorf("expected post content to be 'Last News Content', got %s", result.Content)
			}
			if result.Link != "https://lastnews.com" {
				t.Errorf("expected post link to be 'https://lastnews.com', got %s", result.Link)
			}
			if result.PubTime == "" {
				t.Error("expected PubTime to not be empty")
			}
		},
	)

	t.Run(
		"repository error", func(t *testing.T) {
			expectedError := errors.New("repository error")
			repo := &mockRepositoryForFindLast{
				post: nil,
				err:  expectedError,
			}

			useCase := NewFindLastUseCase(repo)
			ctx := context.Background()

			result, err := useCase.Execute(ctx)

			if err == nil {
				t.Error("expected error, got nil")
			}

			if !errors.Is(err, expectedError) {
				t.Errorf("expected error to wrap repository error, got %v", err)
			}

			if result.ID != 0 {
				t.Errorf("expected empty PostDTO on error, got %+v", result)
			}
		},
	)

	t.Run(
		"no posts found", func(t *testing.T) {
			noPostsError := errors.New("no posts found")
			repo := &mockRepositoryForFindLast{
				post: nil,
				err:  noPostsError,
			}

			useCase := NewFindLastUseCase(repo)
			ctx := context.Background()

			result, err := useCase.Execute(ctx)

			if err == nil {
				t.Error("expected error when no posts found, got nil")
			}

			if !errors.Is(err, noPostsError) {
				t.Errorf("expected error to wrap no posts found error, got %v", err)
			}

			if result.ID != 0 {
				t.Errorf("expected empty PostDTO when no posts found, got %+v", result)
			}
		},
	)

	t.Run(
		"database connection error", func(t *testing.T) {
			dbError := errors.New("database connection failed")
			repo := &mockRepositoryForFindLast{
				post: nil,
				err:  dbError,
			}

			useCase := NewFindLastUseCase(repo)
			ctx := context.Background()

			result, err := useCase.Execute(ctx)

			if err == nil {
				t.Error("expected error when database fails, got nil")
			}

			if !errors.Is(err, dbError) {
				t.Errorf("expected error to wrap database error, got %v", err)
			}

			if result.ID != 0 {
				t.Errorf("expected empty PostDTO on database error, got %+v", result)
			}
		},
	)

	t.Run(
		"context cancellation", func(t *testing.T) {
			postID, _ := dom.NewPostID(1)
			post, _ := dom.NewPost("Test Title", "Test Content", "https://example.com", dom.NewPubTime().Time().Unix())
			post.SetID(postID)

			repo := &mockRepositoryForFindLast{
				post: post,
				err:  nil,
			}

			useCase := NewFindLastUseCase(repo)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := useCase.Execute(ctx)
			_ = err
		},
	)
}
