// Package mapper Переводит сущности в DTO и наоборот.
package mapper

import (
	dom "GoNews/internal/domain/post"
	"fmt"
)

// PostDocument - структура для маппинга новости из Mongo.
type PostDocument struct {
	ID      int32  `bson:"_id,omitempty"`
	Title   string `bson:"title"`
	Content string `bson:"content"`
	PubTime int64  `bson:"pub_time"`
	Link    string `bson:"link"`
}

// MapDocToPost - функция для маппинга новости из Mongo.
func MapDocToPost(doc PostDocument) (*dom.Post, error) {
	id, err := dom.NewPostID(doc.ID)
	if err != nil {
		return nil, fmt.Errorf("MapDocToPost.NewPostID: %w", err)
	}

	title, err := dom.NewPostTitle(doc.Title)
	if err != nil {
		return nil, fmt.Errorf("MapDocToPost.NewTitle: %w", err)
	}

	content, err := dom.NewPostContent(doc.Content)
	if err != nil {
		return nil, fmt.Errorf("MapDocToPost.NewContent: %w", err)
	}

	pubTime, err := dom.NewFromUnixSeconds(doc.PubTime)
	if err != nil {
		return nil, fmt.Errorf("MapDocToPost.NewFromUnixSeconds: %w", err)
	}
	link, err := dom.NewPostLink(doc.Link)
	if err != nil {
		return nil, fmt.Errorf("MapDocToPost.NewPostLink: %w", err)
	}

	return dom.RehydratePost(id, title, content, pubTime, link), nil
}

// FromPostToDoc маппинг доменной модели новости в MongoDB-документ.
func FromPostToDoc(p *dom.Post) *PostDocument {
	return &PostDocument{
		ID:      p.ID().Value(),
		Title:   p.Title().Value(),
		Content: p.Content().Value(),
		PubTime: p.PubTime().Time().Unix(),
		Link:    p.Link().Value(),
	}
}
