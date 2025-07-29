package post

import (
	dom "GoNews/internal/domain/post"
	"context"
	"errors"
	"testing"
	"time"
)

// mockRepository реализует интерфейс dom.Repository для тестирования
type mockStoreRepository struct {
	posts       []*dom.Post
	storeErr    error
	findByIDErr error
}

func (m *mockStoreRepository) Store(ctx context.Context, post *dom.Post) error {
	if m.storeErr != nil {
		return m.storeErr
	}
	m.posts = append(m.posts, post)
	return nil
}

func (m *mockStoreRepository) FindByID(ctx context.Context, id dom.PostID) (*dom.Post, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}

	for _, post := range m.posts {
		if post.ID() == id {
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}

func (m *mockStoreRepository) FindAll(ctx context.Context) ([]*dom.Post, error) {
	return m.posts, nil
}

func (m *mockStoreRepository) FindLast(ctx context.Context) (*dom.Post, error) {
	return nil, nil
}

func (m *mockStoreRepository) FindLatest(ctx context.Context, limit int) ([]*dom.Post, error) {
	return nil, nil
}

// mockParser реализует интерфейс Parser для тестирования
type mockParser struct {
	items []ParsedRSSDTO
	err   error
}

func (m *mockParser) Parse(url string) ([]ParsedRSSDTO, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.items, nil
}

func TestParseAndStoreUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	// Подготавливаем моковые данные RSS
	rssItems := []ParsedRSSDTO{
		{
			Title:   "Test Title 1",
			Content: "Test Content 1",
			Link:    "https://example1.com",
			PubTime: time.Now().Unix(),
		},
		{
			Title:   "Test Title 2",
			Content: "Test Content 2",
			Link:    "https://example2.com",
			PubTime: time.Now().Unix(),
		},
	}

	parser := &mockParser{
		items: rssItems,
		err:   nil,
	}

	repo := &mockStoreRepository{
		posts:       []*dom.Post{},
		storeErr:    nil,
		findByIDErr: errors.New("post not found"), // Симулируем, что посты не найдены
	}

	useCase := NewParseAndStoreUseCase(repo, parser)

	input := ParseAndStoreInputDTO{
		URL: "https://example.com/rss",
	}

	err := useCase.Execute(ctx, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Проверяем, что посты были сохранены
	if len(repo.posts) != 2 {
		t.Fatalf("Expected 2 posts to be stored, got: %d", len(repo.posts))
	}
}

func TestParseAndStoreUseCase_Execute_InvalidInput(t *testing.T) {
	ctx := context.Background()

	parser := &mockParser{}
	repo := &mockRepository{}
	useCase := NewParseAndStoreUseCase(repo, parser)

	// Пустой URL должен вызвать ошибку валидации
	input := ParseAndStoreInputDTO{
		URL: "",
	}

	err := useCase.Execute(ctx, input)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestParseAndStoreUseCase_Execute_ParserError(t *testing.T) {
	ctx := context.Background()

	parser := &mockParser{
		items: nil,
		err:   errors.New("failed to parse RSS"),
	}

	repo := &mockRepository{}
	useCase := NewParseAndStoreUseCase(repo, parser)

	input := ParseAndStoreInputDTO{
		URL: "https://example.com/rss",
	}

	err := useCase.Execute(ctx, input)

	if err == nil {
		t.Fatal("Expected parser error, got nil")
	}

	if !errors.Is(err, parser.err) && err.Error() != "ParseAndStoreUseCase.Parse: failed to parse RSS" {
		t.Fatalf("Expected parser error, got: %v", err)
	}
}
