package comment

import (
	dom "github.com/ee-crocush/go-news/go-comments/internal/domain/comment"
)

// AllByNewsIDDTO представляет входной DTO получения всех комментариев по ID новости.
type AllByNewsIDDTO struct {
	NewsID int32 `json:"news_id"`
}

// CommentDTO представляет выходной DTO коммента.
type CommentDTO struct {
	ID       int64        `json:"id"`
	NewsID   int32        `json:"news_id"`
	ParentID *int64       `json:"parent_id,omitempty"`
	Username string       `json:"username"`
	Content  string       `json:"content"`
	PubTime  string       `json:"pub_time"`
	Children []CommentDTO `json:"children,omitempty"`
}

// MapTreeToDTO переводит дерево сущность в дерево DTO.
func MapTreeToDTO(comments []*dom.Comment) []CommentDTO {
	result := make([]CommentDTO, 0, len(comments))
	for _, c := range comments {
		result = append(result, mapCommentToDTO(c))
	}

	return result
}

// mapCommentToDTO переводит сущность в DTO.
func mapCommentToDTO(comment *dom.Comment) CommentDTO {
	dto := CommentDTO{
		ID:       comment.ID().Value(),
		NewsID:   comment.NewsID().Value(),
		Username: comment.Username().Value(),
		Content:  comment.Content().Value(),
		PubTime:  comment.PubTime().String(),
	}

	if pid := comment.ParentID().Value(); pid != nil {
		dto.ParentID = pid
	}

	for _, child := range comment.Children() {
		dto.Children = append(dto.Children, mapCommentToDTO(child))
	}

	return dto
}
