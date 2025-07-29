// Package post содержит определения бизнес-правил и логики для сущности Post.
package post

import (
	"fmt"
)

// Post представляет новость.
type Post struct {
	id      PostID
	title   PostTitle
	content PostContent
	pubTime PubTime
	link    PostLink
}

// NewPost создает новую новость.
func NewPost(title, content, link string, pubTime int64) (*Post, error) {
	postTitle, err := NewPostTitle(title)
	if err != nil {
		return nil, fmt.Errorf("NewPost.NewPostTitle: %w", err)
	}

	postContent, err := NewPostContent(content)
	if err != nil {
		return nil, fmt.Errorf("NewPost.NewPostTitle: %w", err)
	}

	postLink, err := NewPostLink(link)
	if err != nil {
		return nil, fmt.Errorf("NewPost.NewPostLink: %w", err)
	}

	postPubTime, err := NewFromUnixSeconds(pubTime)
	if err != nil {
		return nil, fmt.Errorf("NewPost.NewPubTime: %w", err)
	}

	return &Post{
		title:   postTitle,
		content: postContent,
		pubTime: postPubTime,
		link:    postLink,
	}, nil
}

// ID возвращает идентификатор новости.
func (p *Post) ID() PostID { return p.id }

// Title возвращает заголовок новости.
func (p *Post) Title() PostTitle { return p.title }

// Content возвращает содержимое новости.
func (p *Post) Content() PostContent { return p.content }

// PubTime возвращает дату публикации новости.
func (p *Post) PubTime() PubTime { return p.pubTime }

// Link возвращает источник новости.
func (p *Post) Link() PostLink { return p.link }

// RehydratePost — вспомогательный конструктор для «восстановления» сущности из БД.
func RehydratePost(id PostID, title PostTitle, content PostContent, pubTime PubTime, link PostLink) *Post {
	return &Post{
		id:      id,
		title:   title,
		content: content,
		pubTime: pubTime,
		link:    link,
	}
}

// SetID устанавливает идентификатор новости.
func (p *Post) SetID(id PostID) { p.id = id }
