// Package usecase выполняет бизнес-логику приложения.
package post

import (
	"fmt"
	dom "github.com/ee-crocush/go-news/go-news/internal/domain/post"
	"net/url"
	"unicode/utf8"
)

// ParsedRSSDTO представляет данные, извлечённые из RSS.
type ParsedRSSDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	PubTime int64  `json:"pub_time"`
}

// ParseAndStoreInputDTO представляет входной DTO для парасинга новостей.
type ParseAndStoreInputDTO struct {
	URL string `json:"url"`
}

func (p *ParseAndStoreInputDTO) Validate() error {
	if p.URL == "" {
		return fmt.Errorf("URL is required")
	}

	_, err := url.ParseRequestURI(p.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	return nil
}

// FindByIDInputDTO представляет входной DTO для поиска поста по ID.
type FindByIDInputDTO struct {
	ID int32 `json:"id"`
}

// FindAllInputDTO представляет входной DTO для поиска поста по параметрам.
type FindAllInputDTO struct {
	Search string
	Limit  int
	Page   int
}

// PostDTO представляет выходной DTO поста.
type PostDTO struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	PubTime string `json:"pub_time"`
}

// FindLatestInputDTO входные данные для поиска последних n новостей.
type FindLatestInputDTO struct {
	Limit int
}

// FindByTitleSubstringInputDTO представляет входной DTO для поиска новостей по заголовку.
type FindByTitleSubstringInputDTO struct {
	Substring string `json:"substring"`
}

// Validate проверяет входные данные для поиска последних n новостей.
func (f *FindLatestInputDTO) Validate() {
	if f.Limit <= 0 {
		f.Limit = 10
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
}

const ContentLimit = 300

// MapPostsToDTO мапит слайс доменных постов в слайс DTO.
func MapPostsToDTO(posts []*dom.Post) []PostDTO {
	postsDTO := make([]PostDTO, 0, len(posts))
	for _, post := range posts {
		postsDTO = append(
			postsDTO, PostDTO{
				ID:      post.ID().Value(),
				Title:   post.Title().Value(),
				Content: trancateContent(post.Content().Value(), ContentLimit),
				Link:    post.Link().Value(),
				PubTime: post.PubTime().String(),
			},
		)
	}
	return postsDTO
}

func trancateContent(content string, limit int) string {
	if utf8.RuneCountInString(content) <= limit {
		return content
	}

	runse := []rune(content)

	return string(runse[:limit])
}
