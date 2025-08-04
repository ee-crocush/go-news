package comment

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewComment(t *testing.T) {
	tests := []struct {
		name     string
		newsID   int32
		username string
		content  string
		pubTime  int64
		wantErr  bool
		errType  error
	}{
		{
			name:     "valid comment with parentID",
			username: "username",
			newsID:   1,
			content:  "Test Content",
			pubTime:  time.Now().Unix(),
			wantErr:  false,
		},
		{
			name:     "valid comment without parentID",
			username: "username",
			newsID:   1,
			content:  "Test Content",
			pubTime:  time.Now().Unix(),
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			newsID:   1,
			content:  "Test Content",
			pubTime:  time.Now().Unix(),
			wantErr:  true,
			errType:  ErrWrongLengthUserName,
		},
		{
			name:     "empty content",
			username: "username",
			newsID:   1,
			content:  "",
			pubTime:  time.Now().Unix(),
			wantErr:  true,
			errType:  ErrEmptyContent,
		},
		{
			name:     "empty news id",
			username: "username",
			newsID:   0,
			content:  "dasdada",
			pubTime:  time.Now().Unix(),
			wantErr:  true,
			errType:  ErrInvalidNewsID,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				comment, err := NewComment(tt.newsID, tt.username, tt.content)

				if tt.wantErr {
					if err == nil {
						t.Errorf("NewComment() expected error but got none")
						return
					}
					if tt.errType != nil && !isErrorInChain(err, tt.errType) {
						t.Errorf("NewComment() expected error type %v in chain, got %v", tt.errType, err)
					}
					return
				}

				if err != nil {
					t.Errorf("NewComment() unexpected error: %v", err)
					return
				}

				if comment == nil {
					t.Errorf("NewComment() returned nil comment")
					return
				}
			},
		)
	}
}

func TestNewComment_ErrorMessages(t *testing.T) {
	// Проверяем, что ошибки правильно оборачиваются
	_, err := NewComment(1, "", "content")
	if !errors.Is(err, ErrWrongLengthUserName) {
		t.Errorf("expected ErrWrongLengthUserName, got %v", err)
	}

	_, err = NewComment(1, "username", "")
	if !errors.Is(err, ErrEmptyContent) {
		t.Errorf("expected ErrEmptyContent, got %v", err)
	}

	_, err = NewComment(0, "username", "")
	if !errors.Is(err, ErrInvalidNewsID) {
		t.Errorf("expected ErrInvalidNewsID, got %v", err)
	}
}

func TestComment_Getters(t *testing.T) {
	// Создаем тестовый коммент
	newsId := int32(1)
	username := "username"
	content := "Test Content"
	pubTime := time.Now().Unix()

	comment, err := NewComment(newsId, username, content)
	if err != nil {
		t.Fatalf("Failed to create test comment: %v", err)
	}

	// Проверяем геттеры
	if comment.Username().Value() != username {
		t.Errorf("Title() = %v, want %v", comment.Username().Value(), username)
	}
	if comment.Content().Value() != content {
		t.Errorf("Content() = %v, want %v", comment.Content().Value(), content)
	}
	if comment.PubTime().Time().Unix() != pubTime {
		t.Errorf("PubTime() = %v, want %v", comment.PubTime().Time().Unix(), pubTime)
	}

	// ID и parentID должен быть нулевым для нового коммента
	if comment.ID().Value() != 0 {
		t.Errorf("ID() = %v, want 0 for new comment", comment.ID().Value())
	}
	if comment.ParentID().Value() != nil {
		t.Errorf("ID() = %v, want nil for new comment", comment.ParentID().Value())
	}
}

func TestComment_SetID(t *testing.T) {
	comment, err := NewComment(1, "username", "Content")
	if err != nil {
		t.Fatalf("Failed to create test comment: %v", err)
	}

	// Тестируем установку валидного ID
	id, err := NewID(123)
	if err != nil {
		t.Fatalf("Failed to create test ID: %v", err)
	}

	comment.SetID(id)

	if !comment.ID().Equal(id) {
		t.Errorf("SetID() failed, got %v, want %v", comment.ID(), id)
	}

	// Тестируем parentID
	newParentID, err := NewParentID(456)
	if err != nil {
		t.Fatalf("Failed to create new test parentID: %v", err)
	}

	comment.SetParentID(newParentID)

	if *comment.ParentID().Value() == 0 {
		t.Errorf("ID() = %v, want %v for new comment", comment.ParentID().Value(), newParentID)
	}
}

func TestRehydrateComment(t *testing.T) {
	// Создаем все необходимые value objects
	id, err := NewID(1)
	if err != nil {
		t.Fatalf("Failed to create CommentID: %v", err)
	}

	NewsId, err := NewNewsID(1)
	if err != nil {
		t.Fatalf("Failed to create NewsId: %v", err)
	}

	parentID, err := NewParentID(1)
	if err != nil {
		t.Fatalf("Failed to create parentID: %v", err)
	}

	username, err := NewUserName("Username")
	if err != nil {
		t.Fatalf("Failed to create Username: %v", err)
	}

	content, err := NewContent("Test Content")
	if err != nil {
		t.Fatalf("Failed to create Content: %v", err)
	}

	testTime := time.Now().Unix()
	pubTime, err := NewFromUnixSeconds(testTime)
	if err != nil {
		t.Fatalf("Failed to create PubTime: %v", err)
	}

	// Тестируем RehydrateComment
	comment := RehydrateComment(id, NewsId, parentID, username, content, pubTime)

	if comment == nil {
		t.Fatal("RehydrateComment() returned nil")
	}

	// Проверяем все поля
	if !comment.ID().Equal(id) {
		t.Errorf("ID() = %v, want %v", comment.ID(), id)
	}
	if comment.ParentID().Value() != parentID.Value() {
		t.Errorf("Title() = %v, want %v", comment.ParentID().Value(), parentID.Value())
	}
	if comment.Username().Value() != username.Value() {
		t.Errorf("Content() = %v, want %v", comment.Username().Value(), username.Value())
	}
	if comment.Content().Value() != content.Value() {
		t.Errorf("Content() = %v, want %v", comment.Content().Value(), content.Value())
	}
	if !comment.PubTime().Time().Equal(pubTime.Time()) {
		t.Errorf("PubTime() = %v, want %v", comment.PubTime().Time(), pubTime.Time())
	}
}

func TestRehydrateComment_WithZeroValues(t *testing.T) {
	// Тестируем RehydrateComment с нулевыми значениями
	var (
		id       ID
		newsID   NewsID
		parentID ParentID
		username UserName
		content  Content
		pubTime  PubTime
	)

	comment := RehydrateComment(id, newsID, parentID, username, content, pubTime)

	if comment == nil {
		t.Fatal("RehydrateComment() returned nil")
	}

	// Проверяем, что все поля установлены (даже если они нулевые)
	if comment.ID().Value() != 0 {
		t.Errorf("ID() = %v, want 0", comment.ID().Value())
	}
	if comment.NewsID().Value() != 0 {
		t.Errorf("ID() = %v, want 0", comment.NewsID().Value())
	}
	if comment.ParentID().Value() != nil {
		t.Errorf("ID() = %v, want nil", comment.ParentID().Value())
	}
	if comment.Username().Value() != "" {
		t.Errorf("Title() = %v, want empty string", comment.Username().Value())
	}
	if comment.Content().Value() != "" {
		t.Errorf("Content() = %v, want empty string", comment.Content().Value())
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

//func ptr[T any](v T) *T {
//	return &v
//}
