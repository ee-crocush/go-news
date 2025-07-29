package handler

import (
	uc "GoNews/internal/usecase/post"
	"GoNews/pkg/api"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// FindByIdRequest - входные данные из тела запроса для получения новости по ID.
type FindByIdRequest struct {
	ID int32 `json:"id"`
}

// FindByIdRequestResponse представляет выходной DTO поста.
type FindByIdRequestResponse struct {
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

	resp := FindByIdRequestResponse{
		ID:      out.ID,
		Title:   out.Title,
		Content: out.Content,
		Link:    out.Link,
		PubTime: out.PubTime,
	}

	return c.Status(fiber.StatusOK).JSON(api.Resp(resp))
}
