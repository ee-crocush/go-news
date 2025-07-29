package post

import (
	"strings"
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
		link    string
		pubTime int64
		wantErr bool
		errType error
	}{
		{
			name:    "valid post",
			title:   "Test Title",
			content: "Test Content",
			link:    "https://example.com",
			pubTime: time.Now().Unix(),
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			content: "Test Content",
			link:    "https://example.com",
			pubTime: time.Now().Unix(),
			wantErr: true,
			errType: ErrEmptyPostTitle,
		},
		{
			name:    "empty content",
			title:   "Test Title",
			content: "",
			link:    "https://example.com",
			pubTime: time.Now().Unix(),
			wantErr: true,
			errType: ErrEmptyPostContent,
		},
		{
			name:    "empty link",
			title:   "Test Title",
			content: "Test Content",
			link:    "",
			pubTime: time.Now().Unix(),
			wantErr: true,
			errType: ErrEmptyPostLink,
		},
		{
			name:    "invalid pub time - zero",
			title:   "Test Title",
			content: "Test Content",
			link:    "https://example.com",
			pubTime: 0,
			wantErr: true,
			errType: ErrEmptyPubTime,
		},
		{
			name:    "invalid pub time - negative",
			title:   "Test Title",
			content: "Test Content",
			link:    "https://example.com",
			pubTime: -1,
			wantErr: true,
			errType: ErrEmptyPubTime,
		},
		{
			name:    "minimal valid data",
			title:   "A",
			content: "B",
			link:    "C",
			pubTime: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				post, err := NewPost(tt.title, tt.content, tt.link, tt.pubTime)

				if tt.wantErr {
					if err == nil {
						t.Errorf("NewPost() expected error but got none")
						return
					}
					if tt.errType != nil && !isErrorInChain(err, tt.errType) {
						t.Errorf("NewPost() expected error type %v in chain, got %v", tt.errType, err)
					}
					return
				}

				if err != nil {
					t.Errorf("NewPost() unexpected error: %v", err)
					return
				}

				if post == nil {
					t.Errorf("NewPost() returned nil post")
					return
				}

				// Проверяем, что поля установлены корректно
				if post.Title().Value() != tt.title {
					t.Errorf("Post.Title() = %v, want %v", post.Title().Value(), tt.title)
				}
				if post.Content().Value() != tt.content {
					t.Errorf("Post.Content() = %v, want %v", post.Content().Value(), tt.content)
				}
				if post.Link().Value() != tt.link {
					t.Errorf("Post.Link() = %v, want %v", post.Link().Value(), tt.link)
				}
				if post.PubTime().Time().Unix() != tt.pubTime {
					t.Errorf("Post.PubTime() = %v, want %v", post.PubTime().Time().Unix(), tt.pubTime)
				}
			},
		)
	}
}

func TestNewPost_ErrorMessages(t *testing.T) {
	// Проверяем, что ошибки правильно оборачиваются
	_, err := NewPost("", "content", "link", time.Now().Unix())
	if err == nil {
		t.Fatal("Expected error for empty title")
	}
	if !strings.Contains(err.Error(), "NewPost.NewPostTitle") {
		t.Errorf("Error should contain context, got: %v", err.Error())
	}

	_, err = NewPost("title", "", "link", time.Now().Unix())
	if err == nil {
		t.Fatal("Expected error for empty content")
	}

	if !strings.Contains(err.Error(), "NewPost.NewPostTitle") {
		t.Errorf("Error should contain context, got: %v", err.Error())
	}

	_, err = NewPost("title", "content", "", time.Now().Unix())
	if err == nil {
		t.Fatal("Expected error for empty link")
	}
	if !strings.Contains(err.Error(), "NewPost.NewPostLink") {
		t.Errorf("Error should contain context, got: %v", err.Error())
	}

	_, err = NewPost("title", "content", "link", 0)
	if err == nil {
		t.Fatal("Expected error for invalid time")
	}
	if !strings.Contains(err.Error(), "NewPost.NewPubTime") {
		t.Errorf("Error should contain context, got: %v", err.Error())
	}
}

func TestPost_Getters(t *testing.T) {
	// Создаем тестовый пост
	title := "Test Title"
	content := "Test Content"
	link := "https://example.com"
	pubTime := time.Now().Unix()

	post, err := NewPost(title, content, link, pubTime)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Проверяем геттеры
	if post.Title().Value() != title {
		t.Errorf("Title() = %v, want %v", post.Title().Value(), title)
	}
	if post.Content().Value() != content {
		t.Errorf("Content() = %v, want %v", post.Content().Value(), content)
	}
	if post.Link().Value() != link {
		t.Errorf("Link() = %v, want %v", post.Link().Value(), link)
	}
	if post.PubTime().Time().Unix() != pubTime {
		t.Errorf("PubTime() = %v, want %v", post.PubTime().Time().Unix(), pubTime)
	}

	// ID должен быть нулевым для нового поста
	if post.ID().Value() != 0 {
		t.Errorf("ID() = %v, want 0 for new post", post.ID().Value())
	}
}

func TestPost_SetID(t *testing.T) {
	post, err := NewPost("Title", "Content", "https://example.com", time.Now().Unix())
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Тестируем установку валидного ID
	id, err := NewPostID(123)
	if err != nil {
		t.Fatalf("Failed to create test ID: %v", err)
	}

	post.SetID(id)

	if !post.ID().Equal(id) {
		t.Errorf("SetID() failed, got %v, want %v", post.ID(), id)
	}

	// Тестируем перезапись ID
	newID, err := NewPostID(456)
	if err != nil {
		t.Fatalf("Failed to create new test ID: %v", err)
	}

	post.SetID(newID)

	if !post.ID().Equal(newID) {
		t.Errorf("SetID() failed to overwrite, got %v, want %v", post.ID(), newID)
	}
}

func TestRehydratePost(t *testing.T) {
	// Создаем все необходимые value objects
	id, err := NewPostID(1)
	if err != nil {
		t.Fatalf("Failed to create PostID: %v", err)
	}

	title, err := NewPostTitle("Test Title")
	if err != nil {
		t.Fatalf("Failed to create PostTitle: %v", err)
	}

	content, err := NewPostContent("Test Content")
	if err != nil {
		t.Fatalf("Failed to create PostContent: %v", err)
	}

	testTime := time.Now().Unix()
	pubTime, err := NewFromUnixSeconds(testTime)
	if err != nil {
		t.Fatalf("Failed to create PubTime: %v", err)
	}

	link, err := NewPostLink("https://example.com")
	if err != nil {
		t.Fatalf("Failed to create PostLink: %v", err)
	}

	// Тестируем RehydratePost
	post := RehydratePost(id, title, content, pubTime, link)

	if post == nil {
		t.Fatal("RehydratePost() returned nil")
	}

	// Проверяем все поля
	if !post.ID().Equal(id) {
		t.Errorf("ID() = %v, want %v", post.ID(), id)
	}
	if post.Title().Value() != title.Value() {
		t.Errorf("Title() = %v, want %v", post.Title().Value(), title.Value())
	}
	if post.Content().Value() != content.Value() {
		t.Errorf("Content() = %v, want %v", post.Content().Value(), content.Value())
	}
	if !post.PubTime().Time().Equal(pubTime.Time()) {
		t.Errorf("PubTime() = %v, want %v", post.PubTime().Time(), pubTime.Time())
	}
	if post.Link().Value() != link.Value() {
		t.Errorf("Link() = %v, want %v", post.Link().Value(), link.Value())
	}
}

func TestRehydratePost_WithZeroValues(t *testing.T) {
	// Тестируем RehydratePost с нулевыми значениями
	var (
		id      PostID
		title   PostTitle
		content PostContent
		pubTime PubTime
		link    PostLink
	)

	post := RehydratePost(id, title, content, pubTime, link)

	if post == nil {
		t.Fatal("RehydratePost() returned nil")
	}

	// Проверяем, что все поля установлены (даже если они нулевые)
	if post.ID().Value() != 0 {
		t.Errorf("ID() = %v, want 0", post.ID().Value())
	}
	if post.Title().Value() != "" {
		t.Errorf("Title() = %v, want empty string", post.Title().Value())
	}
	if post.Content().Value() != "" {
		t.Errorf("Content() = %v, want empty string", post.Content().Value())
	}
	if post.Link().Value() != "" {
		t.Errorf("Link() = %v, want empty string", post.Link().Value())
	}
}

// isErrorInChain проверяет, содержит ли цепочка ошибок указанный тип
func isErrorInChain(err, target error) bool {
	for err != nil {
		if err == target {
			return true
		}
		// Проверяем wrapped errors через fmt.Errorf
		if strings.Contains(err.Error(), target.Error()) {
			return true
		}
		// Проверяем interface для unwrapping
		if unwrapped, ok := err.(interface{ Unwrap() error }); ok {
			err = unwrapped.Unwrap()
		} else {
			break
		}
	}
	return false
}
