package handler

import (
	uc "github.com/ee-crocush/go-news/go-comments/internal/domain/comment"
	"github.com/ee-crocush/go-news/pkg/api"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// TODO: 2025-08-04 Mogush E.E.: Закончить роуты

// CreateRequest - входные данные из тела запроса для создания комментария.
type CreateRequest struct {
	ID int32 `json:"id"`
}

// CreateRequestResponse представляет выходной DTO поста.
type CreateRequestResponse struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	PubTime string `json:"pub_time"`
}

// FindByIDHandler обрабатывает запрос (GET /news/<id>).
func (h *Handler) FindByIDHandler(c *fiber.Ctx) error {
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

	in := uc.FindByIDInputDTO{ID: int32(id)}
	out, err := h.findByIDUC.Execute(c.Context(), in)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Err(err))
	}

	resp := CreateRequestResponse{
		ID:      out.ID,
		Title:   out.Title,
		Content: out.Content,
		Link:    out.Link,
		PubTime: out.PubTime,
	}

	return c.Status(fiber.StatusOK).JSON(api.Resp(resp))
}
