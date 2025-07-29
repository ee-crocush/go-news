package post

import (
	"testing"
	"time"
)

func TestNewPostID(t *testing.T) {
	t.Run(
		"valid ID", func(t *testing.T) {
			id, err := NewPostID(1)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if id.Value() != 1 {
				t.Errorf("expected value 1, got %d", id.Value())
			}
		},
	)

	t.Run(
		"invalid ID - zero", func(t *testing.T) {
			_, err := NewPostID(0)
			if err != ErrInvalidPostID {
				t.Errorf("expected ErrInvalidPostID, got %v", err)
			}
		},
	)

	t.Run(
		"invalid ID - negative", func(t *testing.T) {
			_, err := NewPostID(-1)
			if err != ErrInvalidPostID {
				t.Errorf("expected ErrInvalidPostID, got %v", err)
			}
		},
	)
}

func TestPostID_Equal(t *testing.T) {
	id1, _ := NewPostID(1)
	id2, _ := NewPostID(1)
	id3, _ := NewPostID(2)

	if !id1.Equal(id2) {
		t.Error("expected equal IDs to be equal")
	}

	if id1.Equal(id3) {
		t.Error("expected different IDs to not be equal")
	}
}

func TestNewPostTitle(t *testing.T) {
	t.Run(
		"valid title", func(t *testing.T) {
			title, err := NewPostTitle("Test Title")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if title.Value() != "Test Title" {
				t.Errorf("expected 'Test Title', got %s", title.Value())
			}
		},
	)

	t.Run(
		"empty title", func(t *testing.T) {
			_, err := NewPostTitle("")
			if err != ErrEmptyPostTitle {
				t.Errorf("expected ErrEmptyPostTitle, got %v", err)
			}
		},
	)
}

func TestNewPostContent(t *testing.T) {
	t.Run(
		"valid content", func(t *testing.T) {
			content, err := NewPostContent("Test content")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if content.Value() != "Test content" {
				t.Errorf("expected 'Test content', got %s", content.Value())
			}
		},
	)

	t.Run(
		"empty content", func(t *testing.T) {
			_, err := NewPostContent("")
			if err != ErrEmptyPostContent {
				t.Errorf("expected ErrEmptyPostContent, got %v", err)
			}
		},
	)
}

func TestNewPubTime(t *testing.T) {
	before := time.Now().UTC()
	pubTime := NewPubTime()
	after := time.Now().UTC()

	if pubTime.Time().Before(before) || pubTime.Time().After(after) {
		t.Error("expected PubTime to be between before and after timestamps")
	}
}

func TestNewFromUnixSeconds(t *testing.T) {
	t.Run(
		"valid unix seconds", func(t *testing.T) {
			unixTime := int64(1609459200) // 2021-01-01 00:00:00 UTC
			pubTime, err := NewFromUnixSeconds(unixTime)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			expected := time.Unix(unixTime, 0)
			if !pubTime.Time().Equal(expected) {
				t.Errorf("expected %v, got %v", expected, pubTime.Time())
			}
		},
	)

	t.Run(
		"invalid unix seconds - zero", func(t *testing.T) {
			_, err := NewFromUnixSeconds(0)
			if err != ErrEmptyPubTime {
				t.Errorf("expected ErrEmptyPubTime, got %v", err)
			}
		},
	)

	t.Run(
		"invalid unix seconds - negative", func(t *testing.T) {
			_, err := NewFromUnixSeconds(-1)
			if err != ErrEmptyPubTime {
				t.Errorf("expected ErrEmptyPubTime, got %v", err)
			}
		},
	)
}

func TestPubTime_String(t *testing.T) {
	unixTime := int64(1609459200)
	pubTime, _ := NewFromUnixSeconds(unixTime)

	expected := "2021-01-01 05:00:00"
	if pubTime.String() != expected {
		t.Errorf("expected %s, got %s", expected, pubTime.String())
	}
}

func TestNewPostLink(t *testing.T) {
	t.Run(
		"valid link", func(t *testing.T) {
			link, err := NewPostLink("https://example.com")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if link.Value() != "https://example.com" {
				t.Errorf("expected 'https://example.com', got %s", link.Value())
			}
		},
	)

	t.Run(
		"empty link", func(t *testing.T) {
			_, err := NewPostLink("")
			if err != ErrEmptyPostLink {
				t.Errorf("expected ErrEmptyPostLink, got %v", err)
			}
		},
	)
}
