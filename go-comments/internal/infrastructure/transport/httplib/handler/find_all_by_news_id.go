package handler

import (
	uc "github.com/ee-crocush/go-news/go-comments/internal/usecase/comment"
	"github.com/ee-crocush/go-news/pkg/api"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// FindAllByNewsIDResponse представляет ответ на запрос получения всех комментариев конкретной новости.
// Создавать еще отдельно структуру нет смысла, т.к. придется снова рекурсивно проходится по массиву.
type FindAllByNewsIDResponse struct {
	Comments []uc.CommentDTO `json:"comments"`
}

// FindAllByNewsIDHandler обрабатывает запрос на получение всех комментариев конкретного поста (GET /comments).
func (h *Handler) FindAllByNewsIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(api.ErrWithCode("missing-id", "missing post ID in URL"))
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).
			JSON(api.ErrWithCode("invalid-id", "post ID must be positive integer"))
	}
	in := uc.AllByNewsIDDTO{NewsID: int32(id)}
	out, err := h.findAllByNewsUC.Execute(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Err(err))
	}

	response := FindAllByNewsIDResponse{
		Comments: out,
	}

	return c.Status(fiber.StatusOK).JSON(api.Resp(response))
}
