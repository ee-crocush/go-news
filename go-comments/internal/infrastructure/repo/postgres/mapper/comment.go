// Package mapper Переводит сущности в DTO и наоборот.
package mapper

import (
	"fmt"
	dom "github.com/ee-crocush/go-news/go-comments/internal/domain/comment"
)

// CommentRow - структура для маппинга комментария из PostgreSQL.
type CommentRow struct {
	ID       int64  `json:"id"`
	ParentID *int64 `json:"parent_id,omitempty"`
	Username string `json:"username"`
	Content  string `json:"content"`
	PubTime  int64  `json:"pub_time"`
}

// MapRowToComment - функция для маппинга комментария из PostgreSQL.
func MapRowToComment(row CommentRow) (*dom.Comment, error) {
	id, err := dom.NewID(row.ID)
	if err != nil {
		return nil, fmt.Errorf("MapRowToComment.NewID: %w", err)
	}

	username, err := dom.NewUserName(row.Username)
	if err != nil {
		return nil, fmt.Errorf("MapRowToComment.NewUserName: %w", err)
	}

	content, err := dom.NewContent(row.Content)
	if err != nil {
		return nil, fmt.Errorf("MapRowToComment.NewContent: %w", err)
	}

	pubTime, err := dom.NewFromUnixSeconds(row.PubTime)
	if err != nil {
		return nil, fmt.Errorf("MapRowToComment.NewFromUnixSeconds: %w", err)
	}

	var parentID dom.ParentID
	if row.ParentID == nil {
		parentID = dom.NewEmptyParentID()
	} else {
		parentID, err = dom.NewParentID(*row.ParentID)
		if err != nil {
			return nil, fmt.Errorf("MapRowToComment.NewParentID: %w", err)
		}
	}

	return dom.RehydrateComment(id, parentID, username, content, pubTime), nil
}
