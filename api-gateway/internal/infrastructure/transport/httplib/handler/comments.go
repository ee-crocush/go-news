package handler

import (
	"github.com/gofiber/fiber/v2"
)

// CreateComments Создает новый комментарий.
// @Summary Создать новый комментарий
// @Description Создать новый комментарий для конкретной новости.
// @Tags comments
// @Accept json
// @Produce json
// @Param request body CreateCommentRequest true "Данные нового комментария"
// @Success 200 {object} domain.CommentResponse
// @Router /api/comments [post]
func (h *Handler) CreateComments(c *fiber.Ctx) error {
	return h.handleServiceRequest(
		c, ServiceRequest{
			RouteName: CommentsRouteName,
			Path:      "/comments",
		},
	)
}
