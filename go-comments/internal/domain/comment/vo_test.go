package comment

import (
	"errors"
	"testing"
	"time"
)

func TestNewCommentID(t *testing.T) {
	t.Run(
		"valid ID", func(t *testing.T) {
			id, err := NewID(1)
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
			_, err := NewID(0)
			if !errors.Is(err, ErrInvalidCommentID) {
				t.Errorf("expected ErrInvalidCommentID, got %v", err)
			}
		},
	)

	t.Run(
		"invalid ID - negative", func(t *testing.T) {
			_, err := NewID(-1)
			if !errors.Is(err, ErrInvalidCommentID) {
				t.Errorf("expected ErrInvalidID, got %v", err)
			}
		},
	)
}

func TestID_Equal(t *testing.T) {
	id1, _ := NewID(1)
	id2, _ := NewID(1)
	id3, _ := NewID(2)

	if !id1.Equal(id2) {
		t.Error("expected equal IDs to be equal")
	}

	if id1.Equal(id3) {
		t.Error("expected different IDs to not be equal")
	}
}

func TestNewTestNewParentID(t *testing.T) {
	t.Run(
		"valid parent ID", func(t *testing.T) {
			id, err := NewParentID(1)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if *id.Value() != 1 {
				t.Errorf("expected value 1, got %d", id.Value())
			}
		},
	)

	t.Run(
		"invalid parent ID - zero", func(t *testing.T) {
			_, err := NewParentID(0)
			if !errors.Is(err, ErrInvalidParentID) {
				t.Errorf("expected ErrInvalidParentID, got %v", err)
			}
		},
	)

	t.Run(
		"invalid ID - negative", func(t *testing.T) {
			_, err := NewParentID(-1)
			if !errors.Is(err, ErrInvalidParentID) {
				t.Errorf("expected ErrInvalidParentID, got %v", err)
			}
		},
	)
}

func TestNewContent(t *testing.T) {
	t.Run(
		"valid content", func(t *testing.T) {
			content, err := NewContent("Test content")
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
			_, err := NewContent("")
			if !errors.Is(err, ErrEmptyContent) {
				t.Errorf("expected ErrEmptyContent, got %v", err)
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
			if !errors.Is(err, ErrEmptyPubTime) {
				t.Errorf("expected ErrEmptyPubTime, got %v", err)
			}
		},
	)

	t.Run(
		"invalid unix seconds - negative", func(t *testing.T) {
			_, err := NewFromUnixSeconds(-1)
			if !errors.Is(err, ErrEmptyPubTime) {
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

func TestNewUserName(t *testing.T) {
	t.Run(
		"valid username", func(t *testing.T) {
			content, err := NewUserName("Test username")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if content.Value() != "Test username" {
				t.Errorf("expected 'Test content', got %s", content.Value())
			}
		},
	)

	t.Run(
		"empty username", func(t *testing.T) {
			_, err := NewUserName("")
			if !errors.Is(err, ErrWrongLengthUserName) {
				t.Errorf("expected ErrWrongLengthUserName, got %v", err)
			}
		},
	)

	t.Run(
		"length username < 6", func(t *testing.T) {
			_, err := NewUserName("user")
			if !errors.Is(err, ErrWrongLengthUserName) {
				t.Errorf("expected ErrWrongLengthUserName, got %v", err)
			}
		},
	)

	t.Run(
		"length username > 50", func(t *testing.T) {
			_, err := NewUserName("username12username12username12username12username1234")
			if !errors.Is(err, ErrWrongLengthUserName) {
				t.Errorf("expected ErrWrongLengthUserName, got %v", err)
			}
		},
	)
}
