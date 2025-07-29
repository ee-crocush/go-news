package rss

import (
	"testing"
	"time"
)

func TestParser_Parse(t *testing.T) {
	parser := NewParser(10 * time.Second)

	feed, err := parser.Parse("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed) == 0 {
		t.Fatal("данные не раскодированы")
	}
	t.Logf("получено %d новостей\n", len(feed))
}

func TestParser_parseTime(t *testing.T) {
	parser := NewParser(5 * time.Second)

	testCases := []struct {
		input string
		name  string
	}{
		{
			input: "Mon, 02 Jan 2006 15:04:05 GMT",
			name:  "Standard RFC format",
		},
		{
			input: "2006-01-02T15:04:05Z",
			name:  "ISO format",
		},
		{
			input: "invalid date",
			name:  "Invalid date fallback",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				result := parser.parseTime(tc.input)
				if result <= 0 {
					t.Errorf("Expected valid timestamp, got: %d", result)
				}
				t.Logf("Input: %s, Result: %d", tc.input, result)
			},
		)
	}
}

func TestNewParser(t *testing.T) {
	timeout := 10 * time.Second
	parser := NewParser(timeout)

	if parser == nil {
		t.Fatal("Expected parser instance, got nil")
	}

	if parser.client == nil {
		t.Fatal("Expected HTTP client to be initialized")
	}

	if parser.client.Timeout != timeout {
		t.Errorf("Expected timeout %v, got: %v", timeout, parser.client.Timeout)
	}
}
