package post

import (
	dom "GoNews/internal/domain/post"
	"context"
	"errors"
	"testing"
)

// mockRepositoryForFindByID implements dom.Repository for testing
type mockRepositoryForFindByID struct {
	post *dom.Post
	err  error
}

func (m *mockRepositoryForFindByID) FindByID(ctx context.Context, postID dom.PostID) (*dom.Post, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.post, nil
}

func (m *mockRepositoryForFindByID) FindAll(ctx context.Context) ([]*dom.Post, error) {
	return nil, nil
}

func (m *mockRepositoryForFindByID) FindLast(ctx context.Context) (*dom.Post, error) {
	return nil, nil
}

func (m *mockRepositoryForFindByID) FindLatest(ctx context.Context, limit int) ([]*dom.Post, error) {
	return nil, nil
}

func (m *mockRepositoryForFindByID) Store(ctx context.Context, post *dom.Post) error {
	return nil
}

func TestNewFindByIDUseCase(t *testing.T) {
	repo := &mockRepositoryForFindByID{}
	useCase := NewFindByIDUseCase(repo)

	if useCase == nil {
		t.Error("expected useCase to not be nil")
	}

	if useCase.repo != repo {
		t.Error("expected repo to be set correctly")
	}
}

func TestFindByIDUseCase_Execute(t *testing.T) {
	t.Run(
		"successful execution", func(t *testing.T) {
			// Create mock post
			postID, _ := dom.NewPostID(1)
			post, _ := dom.NewPost("Test Title", "Test Content", "https://example.com", dom.NewPubTime().Time().Unix())
			post.SetID(postID)

			repo := &mockRepositoryForFindByID{
				post: post,
				err:  nil,
			}

			useCase := NewFindByIDUseCase(repo)
			ctx := context.Background()

			input := FindByIDInputDTO{ID: 1}
			result, err := useCase.Execute(ctx, input)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if result.ID != 1 {
				t.Errorf("expected post ID to be 1, got %d", result.ID)
			}
			if result.Title != "Test Title" {
				t.Errorf("expected post title to be 'Test Title', got %s", result.Title)
			}
			if result.Content != "Test Content" {
				t.Errorf("expected post content to be 'Test Content', got %s", result.Content)
			}
			if result.Link != "https://example.com" {
				t.Errorf("expected post link to be 'https://example.com', got %s", result.Link)
			}
			if result.PubTime == "" {
				t.Error("expected PubTime to not be empty")
			}
		},
	)

	t.Run(
		"invalid post ID", func(t *testing.T) {
			repo := &mockRepositoryForFindByID{}
			useCase := NewFindByIDUseCase(repo)
			ctx := context.Background()

			input := FindByIDInputDTO{ID: 0}
			result, err := useCase.Execute(ctx, input)

			if err == nil {
				t.Error("expected error for invalid post ID, got nil")
			}

			if result.ID != 0 {
				t.Errorf("expected empty PostDTO on error, got %+v", result)
			}
		},
	)

	t.Run(
		"negative post ID", func(t *testing.T) {
			repo := &mockRepositoryForFindByID{}
			useCase := NewFindByIDUseCase(repo)
			ctx := context.Background()

			input := FindByIDInputDTO{ID: -1} // Invalid ID
			result, err := useCase.Execute(ctx, input)

			if err == nil {
				t.Error("expected error for negative post ID, got nil")
			}

			if result.ID != 0 {
				t.Errorf("expected empty PostDTO on error, got %+v", result)
			}
		},
	)

	t.Run(
		"repository error", func(t *testing.T) {
			expectedError := errors.New("repository error")
			repo := &mockRepositoryForFindByID{
				post: nil,
				err:  expectedError,
			}

			useCase := NewFindByIDUseCase(repo)
			ctx := context.Background()

			input := FindByIDInputDTO{ID: 1}
			result, err := useCase.Execute(ctx, input)

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
		"post not found", func(t *testing.T) {
			notFoundError := errors.New("post not found")
			repo := &mockRepositoryForFindByID{
				post: nil,
				err:  notFoundError,
			}

			useCase := NewFindByIDUseCase(repo)
			ctx := context.Background()

			input := FindByIDInputDTO{ID: 999}
			result, err := useCase.Execute(ctx, input)

			if err == nil {
				t.Error("expected error when post not found, got nil")
			}

			if !errors.Is(err, notFoundError) {
				t.Errorf("expected error to wrap not found error, got %v", err)
			}

			if result.ID != 0 {
				t.Errorf("expected empty PostDTO when post not found, got %+v", result)
			}
		},
	)

	t.Run(
		"context cancellation", func(t *testing.T) {
			postID, _ := dom.NewPostID(1)
			post, _ := dom.NewPost("Test Title", "Test Content", "https://example.com", dom.NewPubTime().Time().Unix())
			post.SetID(postID)

			repo := &mockRepositoryForFindByID{
				post: post,
				err:  nil,
			}

			useCase := NewFindByIDUseCase(repo)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			input := FindByIDInputDTO{ID: 1}
			_, err := useCase.Execute(ctx, input)
			_ = err
		},
	)
}
