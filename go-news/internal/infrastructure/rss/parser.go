// Package rss выполняет парсинг из url с постами в DTO.
package rss

import (
	uc "GoNews/internal/usecase/post"
	"encoding/xml"
	"fmt"
	strip "github.com/grokify/html-strip-tags-go"
	"io"
	"net/http"
	"strings"
	"time"
)

// Feed представляет RSS ленту.
type Feed struct {
	Channel Channel `xml:"channel"`
}

// Channel представляет канал RSS.
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

// Item представляет элемент RSS.
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

// Parser представляет парсер RSS ленты.
type Parser struct {
	client *http.Client
}

// NewParser создает новый экземпляр парсера RSS ленты.
func NewParser(timeout time.Duration) *Parser {
	return &Parser{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Parse парсит RSS ленту по указанному URL и возвращает слайс спарсенных DTO.
func (p *Parser) Parse(url string) ([]uc.ParsedRSSDTO, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http parse get error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http read body error: %w", err)
	}

	var feed Feed
	if err = xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("xml unmarshal error: %w", err)
	}

	var posts []uc.ParsedRSSDTO
	for _, item := range feed.Channel.Items {
		post := p.itemToDTO(item)
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *Parser) itemToDTO(item Item) uc.ParsedRSSDTO {
	content := strip.StripTags(item.Description)
	pubTime := p.parseTime(item.PubDate)

	return uc.ParsedRSSDTO{
		Title:   item.Title,
		Content: content,
		Link:    item.Link,
		PubTime: pubTime,
	}
}

func (p *Parser) parseTime(dateStr string) int64 {
	dateStr = strings.ReplaceAll(dateStr, ",", "")
	formats := []string{
		"Mon 2 Jan 2006 15:04:05 -0700",
		"Mon 2 Jan 2006 15:04:05 GMT",
		"2 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Unix()
		}
	}
	return time.Now().Unix()
}
