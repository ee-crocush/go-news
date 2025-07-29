package post

import (
	dom "GoNews/internal/domain/post"
	"context"
	"errors"
	"testing"
)

// mockRepositoryForFindLatest implements dom.Repository for testing
type mockRepositoryForFindLatest struct {
	posts []*dom.Post
	err   error
}

func (m *mockRepositoryForFindLatest) FindLatest(ctx context.Context, limit int) ([]*dom.Post, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.posts, nil
}

func (m *mockRepositoryForFindLatest) FindByID(ctx context.Context, postID dom.PostID) (*dom.Post, error) {
	return nil, nil
}

func (m *mockRepositoryForFindLatest) FindAll(ctx context.Context) ([]*dom.Post, error) {
	return nil, nil
}

func (m *mockRepositoryForFindLatest) FindLast(ctx context.Context) (*dom.Post, error) {
	return nil, nil
}

func (m *mockRepositoryForFindLatest) Store(ctx context.Context, post *dom.Post) error {
	return nil
}

func TestNewFindLatestUseCase(t *testing.T) {
	repo := &mockRepositoryForFindLatest{}
	useCase := NewFindLatestUseCase(repo)

	if useCase == nil {
		t.Error("expected useCase to not be nil")
	}

	if useCase.repo != repo {
		t.Error("expected repo to be set correctly")
	}
}

func TestFindLatestUseCase_Execute(t *testing.T) {
	t.Run(
		"successful execution with multiple posts", func(t *testing.T) {
			postID1, _ := dom.NewPostID(1)
			post1, _ := dom.NewPost(
				"Latest Title 1", "Latest Content 1", "https://latest1.com", dom.NewPubTime().Time().Unix(),
			)
			post1.SetID(postID1)

			postID2, _ := dom.NewPostID(2)
			post2, _ := dom.NewPost(
				"Latest Title 2", "Latest Content 2", "https://latest2.com", dom.NewPubTime().Time().Unix(),
			)
			post2.SetID(postID2)

			postID3, _ := dom.NewPostID(3)
			post3, _ := dom.NewPost(
				"Latest Title 3", "Latest Content 3", "https://latest3.com", dom.NewPubTime().Time().Unix(),
			)
			post3.SetID(postID3)

			mockPosts := []*dom.Post{post1, post2, post3}

			repo := &mockRepositoryForFindLatest{
				posts: mockPosts,
				err:   nil,
			}

			useCase := NewFindLatestUseCase(repo)
			ctx := context.Background()

			input := FindLatestInputDTO{Limit: 3}
			result, err := useCase.Execute(ctx, input)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(result) != 3 {
				t.Errorf("expected 3 posts, got %d", len(result))
			}

			if result[0].ID != 1 {
				t.Errorf("expected first post ID to be 1, got %d", result[0].ID)
			}
			if result[0].Title != "Latest Title 1" {
				t.Errorf("expected first post title to be 'Latest Title 1', got %s", result[0].Title)
			}

			if result[1].ID != 2 {
				t.Errorf("expected second post ID to be 2, got %d", result[1].ID)
			}
			if result[1].Title != "Latest Title 2" {
				t.Errorf("expected second post title to be 'Latest Title 2', got %s", result[1].Title)
			}
		},
	)

	t.Run(
		"successful execution with limit 1", func(t *testing.T) {
			postID, _ := dom.NewPostID(1)
			post, _ := dom.NewPost(
				"Single Latest Title", "Single Latest Content", "https://single.com", dom.NewPubTime().Time().Unix(),
			)
			post.SetID(postID)

			mockPosts := []*dom.Post{post}

			repo := &mockRepositoryForFindLatest{
				posts: mockPosts,
				err:   nil,
			}

			useCase := NewFindLatestUseCase(repo)
			ctx := context.Background()

			input := FindLatestInputDTO{Limit: 1}
			result, err := useCase.Execute(ctx, input)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(result) != 1 {
				t.Errorf("expected 1 post, got %d", len(result))
			}

			if result[0].ID != 1 {
				t.Errorf("expected post ID to be 1, got %d", result[0].ID)
			}
			if result[0].Title != "Single Latest Title" {
				t.Errorf("expected post title to be 'Single Latest Title', got %s", result[0].Title)
			}
		},
	)

	t.Run(
		"empty result", func(t *testing.T) {
			repo := &mockRepositoryForFindLatest{
				posts: []*dom.Post{},
				err:   nil,
			}

			useCase := NewFindLatestUseCase(repo)
			ctx := context.Background()

			input := FindLatestInputDTO{Limit: 5}
			result, err := useCase.Execute(ctx, input)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(result) != 0 {
				t.Errorf("expected 0 posts, got %d", len(result))
			}
		},
	)

	t.Run(
		"repository error", func(t *testing.T) {
			expectedError := errors.New("repository error")
			repo := &mockRepositoryForFindLatest{
				posts: nil,
				err:   expectedError,
			}

			useCase := NewFindLatestUseCase(repo)
			ctx := context.Background()

			input := FindLatestInputDTO{Limit: 5}
			result, err := useCase.Execute(ctx, input)

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
		"database connection error", func(t *testing.T) {
			dbError := errors.New("database connection failed")
			repo := &mockRepositoryForFindLatest{
				posts: nil,
				err:   dbError,
			}

			useCase := NewFindLatestUseCase(repo)
			ctx := context.Background()

			input := FindLatestInputDTO{Limit: 10}
			result, err := useCase.Execute(ctx, input)

			if err == nil {
				t.Error("expected error when database fails, got nil")
			}

			if !errors.Is(err, dbError) {
				t.Errorf("expected error to wrap database error, got %v", err)
			}

			if len(result) != 0 {
				t.Errorf("expected empty result on database error, got %d posts", len(result))
			}
		},
	)

	t.Run(
		"large limit", func(t *testing.T) {
			// Create mock posts
			postID1, _ := dom.NewPostID(1)
			post1, _ := dom.NewPost("Title 1", "Content 1", "https://example1.com", dom.NewPubTime().Time().Unix())
			post1.SetID(postID1)

			postID2, _ := dom.NewPostID(2)
			post2, _ := dom.NewPost("Title 2", "Content 2", "https://example2.com", dom.NewPubTime().Time().Unix())
			post2.SetID(postID2)

			mockPosts := []*dom.Post{post1, post2}

			repo := &mockRepositoryForFindLatest{
				posts: mockPosts,
				err:   nil,
			}

			useCase := NewFindLatestUseCase(repo)
			ctx := context.Background()

			input := FindLatestInputDTO{Limit: 100} // Large limit
			result, err := useCase.Execute(ctx, input)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(result) != 2 {
				t.Errorf("expected 2 posts (all available), got %d", len(result))
			}
		},
	)

	t.Run(
		"context cancellation", func(t *testing.T) {
			postID, _ := dom.NewPostID(1)
			post, _ := dom.NewPost("Test Title", "Test Content", "https://example.com", dom.NewPubTime().Time().Unix())
			post.SetID(postID)

			repo := &mockRepositoryForFindLatest{
				posts: []*dom.Post{post},
				err:   nil,
			}

			useCase := NewFindLatestUseCase(repo)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			input := FindLatestInputDTO{Limit: 5}
			_, err := useCase.Execute(ctx, input)
			_ = err
		},
	)
}
