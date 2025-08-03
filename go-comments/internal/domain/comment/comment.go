// Package comment содержит определения бизнес-правил и логики для сущности "Комментарий".
package comment

import (
	"fmt"
)

// Comment представляет комментарий.
type Comment struct {
	id       ID
	parentID ParentID
	username UserName
	content  Content
	pubTime  PubTime
}

// NewComment создает новый комментарий.
func NewComment(username, content string) (*Comment, error) {
	usernameVO, err := NewUserName(username)
	if err != nil {
		return nil, fmt.Errorf("NewComment.NewUserName: %w", err)
	}

	contentVO, err := NewContent(content)
	if err != nil {
		return nil, fmt.Errorf("NewPost.NewContent: %w", err)
	}

	return &Comment{
		username: usernameVO,
		content:  contentVO,
		pubTime:  NewPubTime(),
	}, nil
}

// ID возвращает идентификатор новости.
func (c *Comment) ID() ID { return c.id }

// ParentID возвращает идентификатор родительского комментария.
func (c *Comment) ParentID() ParentID { return c.parentID }

// Username возвращает имя пользователя, оставившего комментарий.
func (c *Comment) Username() UserName { return c.username }

// Content возвращает содержимое комментария.
func (c *Comment) Content() Content { return c.content }

// PubTime возвращает время публикации комментария.
func (c *Comment) PubTime() PubTime { return c.pubTime }

// RehydrateComment — вспомогательный конструктор для «восстановления» сущности из БД.
func RehydrateComment(id ID, parentID ParentID, username UserName, content Content, pubTime PubTime) *Comment {
	return &Comment{
		id:       id,
		parentID: parentID,
		username: username,
		content:  content,
		pubTime:  pubTime,
	}
}

// SetID устанавливает идентификатор комментария.
func (c *Comment) SetID(id ID) { c.id = id }

// SetParentID устанавливает идентификатор родительского комментария.
func (c *Comment) SetParentID(id ParentID) { c.parentID = id }
