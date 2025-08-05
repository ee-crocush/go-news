package handler

import (
	"github.com/gofiber/fiber/v2"
)

const CommentsRouteName = "go-comments"

// CreateCommentRequest представляет тело запроса для создания комментария.
type CreateCommentRequest struct {
	NewsID   int32  `json:"news_id" example:"1"`
	ParentID *int64 `json:"parent_id" example:"1"`
	Username string `json:"username" example:"Example_username"`
	Content  string `json:"content" example:"Example content"`
}

// CommentResponse описывает структуру ответа на /comments.
type CommentResponse struct {
	Status  string `json:"status" example:"OK"`
	Message string `json:"message" example:"Example message success"`
}

// FindAllCommentsByNewsID получает все комментарии для конкретной новости.
// @Summary Получить все комментарии по ID новости
// @Description Возвращает список всех комментариев по ID новости.
// @Tags comments
// @Produce json
// @Success 200 {array} CommentResponse
// @Router /api/comments [get]
func (h *Handler) FindAllCommentsByNewsID(c *fiber.Ctx) error {
	return h.proxyRequest(c, CommentsRouteName, "/comments")
}

// CreateComments Создает новый комментарий.
// @Summary Создать новый комментарий
// @Description Создать новый комментарий для конкретной новости.
// @Tags comments
// @Accept json
// @Produce json
// @Param request body CreateCommentRequest true "Данные нового комментария"
// @Success 200 {object} CommentResponse
// @Router /api/comments [post]
func (h *Handler) CreateComments(c *fiber.Ctx) error {
	return h.proxyRequest(c, CommentsRouteName, "/comments")
}
