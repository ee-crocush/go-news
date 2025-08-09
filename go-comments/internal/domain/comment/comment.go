// Package comment содержит определения бизнес-правил и логики для сущности "Комментарий".
package comment

import (
	"fmt"
)

// Comment представляет комментарий.
type Comment struct {
	id        ID
	newsID    NewsID
	parentID  ParentID
	username  UserName
	content   Content
	pubTime   CommentTime
	createdAt CommentTime
	status    Status
	children  []*Comment
}

// NewComment создает новый комментарий.
func NewComment(newsID int32, username, content string) (*Comment, error) {
	newsIDVO, err := NewNewsID(newsID)
	if err != nil {
		return nil, fmt.Errorf("NewComment.NewNewsID: %w", err)
	}

	usernameVO, err := NewUserName(username)
	if err != nil {
		return nil, fmt.Errorf("NewComment.NewUserName: %w", err)
	}

	contentVO, err := NewContent(content)
	if err != nil {
		return nil, fmt.Errorf("NewPost.NewContent: %w", err)
	}

	statusVO, err := NewStatus(Pending)
	if err != nil {
		return nil, fmt.Errorf("NewPost.NewStatus: %w", err)
	}

	return &Comment{
		newsID:    newsIDVO,
		username:  usernameVO,
		content:   contentVO,
		createdAt: NewTime(),
		status:    statusVO,
	}, nil
}

// Геттеры

// ID возвращает идентификатор коммента.
func (c *Comment) ID() ID { return c.id }

// NewsID возвращает идентификатор новости.
func (c *Comment) NewsID() NewsID { return c.newsID }

// ParentID возвращает идентификатор родительского комментария.
func (c *Comment) ParentID() ParentID { return c.parentID }

// Username возвращает имя пользователя, оставившего комментарий.
func (c *Comment) Username() UserName { return c.username }

// Content возвращает содержимое комментария.
func (c *Comment) Content() Content { return c.content }

// CreatedAt возвращает время создания комментария.
func (c *Comment) CreatedAt() CommentTime { return c.createdAt }

// PubTime возвращает время публикации комментария.
func (c *Comment) PubTime() CommentTime { return c.pubTime }

// Status возвращает статус модерации.
func (c *Comment) Status() Status { return c.status }

// Children возвращает дочерние комментарии.
func (c *Comment) Children() []*Comment {
	return c.children
}

// IsApproved возвращает true, если комментарий прошел модерацию.
func (c *Comment) IsApproved() bool {
	return c.status.Value() == Approved
}

// Сеттеры

// AddChild добавляет дочерний комментарий.
func (c *Comment) AddChild(child *Comment) {
	c.children = append(c.children, child)
}

// SetID устанавливает идентификатор комментария.
func (c *Comment) SetID(id ID) { c.id = id }

// SetParentID устанавливает идентификатор родительского комментария.
func (c *Comment) SetParentID(id ParentID) { c.parentID = id }

// SetStatus устанавливает статус комментария.
func (c *Comment) SetStatus(status Status) {
	c.status = status
}

// RehydrateComment — вспомогательный конструктор для «восстановления» сущности из БД.
func RehydrateComment(
	id ID, newsID NewsID, parentID ParentID, username UserName, content Content, pubTime CommentTime, status Status,
) *Comment {
	return &Comment{
		id:       id,
		newsID:   newsID,
		parentID: parentID,
		username: username,
		content:  content,
		pubTime:  pubTime,
		status:   status,
	}
}
