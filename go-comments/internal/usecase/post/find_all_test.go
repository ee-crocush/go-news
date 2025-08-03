package post

import (
	dom "GoNews/internal/domain/post"
	"context"
	"errors"
	"testing"
)

// mockRepository implements Repository for testing
type mockRepository struct {
	posts []*dom.Post
	err   error
}

func (m *mockRepository) FindAll(ctx context.Context) ([]*dom.Post, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.posts, nil
}

func (m *mockRepository) FindByID(ctx context.Context, postID dom.PostID) (*dom.Post, error) {
	return nil, nil
}

func (m *mockRepository) FindLast(ctx context.Context) (*dom.Post, error) {
	return nil, nil
}

func (m *mockRepository) FindLatest(ctx context.Context, limit int) ([]*dom.Post, error) {
	return nil, nil
}

func (m *mockRepository) Store(ctx context.Context, post *dom.Post) error {
	return nil
}

func TestNewFindAllUseCase(t *testing.T) {
	repo := &mockRepository{}
	useCase := NewFindAllUseCase(repo)

	if useCase == nil {
		t.Error("expected useCase to not be nil")
	}

	if useCase.repo != repo {
		t.Error("expected repo to be set correctly")
	}
}

func TestFindAllUseCase_Execute(t *testing.T) {
	t.Run(
		"successful execution", func(t *testing.T) {
			postID1, _ := dom.NewPostID(1)
			post1, _ := dom.NewPost("Title 1", "Content 1", "https://example1.com", dom.NewPubTime().Time().Unix())
			post1.SetID(postID1)

			postID2, _ := dom.NewPostID(2)
			post2, _ := dom.NewPost("Title 2", "Content 2", "https://example2.com", dom.NewPubTime().Time().Unix())
			post2.SetID(postID2)

			mockPosts := []*dom.Post{
				post1,
				post2,
			}

			repo := &mockRepository{
				posts: mockPosts,
				err:   nil,
			}

			useCase := NewFindAllUseCase(repo)
			ctx := context.Background()

			result, err := useCase.Execute(ctx)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(result) != 2 {
				t.Errorf("expected 2 posts, got %d", len(result))
			}

			if result[0].ID != 1 {
				t.Errorf("expected first post ID to be 1, got %d", result[0].ID)
			}
			if result[0].Title != "Title 1" {
				t.Errorf("expected first post title to be 'Title 1', got %s", result[0].Title)
			}

			if result[1].ID != 2 {
				t.Errorf("expected second post ID to be 2, got %d", result[1].ID)
			}
			if result[1].Title != "Title 2" {
				t.Errorf("expected second post title to be 'Title 2', got %s", result[1].Title)
			}
		},
	)

	t.Run(
		"repository error", func(t *testing.T) {
			expectedError := errors.New("repository error")
			repo := &mockRepository{
				posts: nil,
				err:   expectedError,
			}

			useCase := NewFindAllUseCase(repo)
			ctx := context.Background()

			result, err := useCase.Execute(ctx)

			if err == nil {
				t.Error("expected error, got nil")
			}

			if !errors.Is(err, expectedError) {
				t.Errorf("expected error to wrap repository error, got %v", err)
			}

			if len(result) != 0 {
				t.Errorf("expected empty result on error, got %d posts", len(result))
			}
		},
	)

	t.Run(
		"empty repository", func(t *testing.T) {
			repo := &mockRepository{
				posts: []*dom.Post{},
				err:   nil,
			}

			useCase := NewFindAllUseCase(repo)
			ctx := context.Background()

			result, err := useCase.Execute(ctx)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(result) != 0 {
				t.Errorf("expected 0 posts, got %d", len(result))
			}
		},
	)

	t.Run(
		"context cancellation", func(t *testing.T) {
			repo := &mockRepository{
				posts: []*dom.Post{},
				err:   nil,
			}

			useCase := NewFindAllUseCase(repo)
			ctx, cancel := context.WithCancel(context.Background())
			cancel() //
			_, err := useCase.Execute(ctx)
			_ = err
		},
	)
}
