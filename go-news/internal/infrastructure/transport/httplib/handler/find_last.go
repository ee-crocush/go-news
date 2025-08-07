package handler

import (
	"github.com/ee-crocush/go-news/pkg/api"
	"github.com/gofiber/fiber/v2"
)

// FindLastResponse представляет ответ на запрос получения последней новости.
type FindLastResponse struct {
	Post PostDTO `json:"post"`
}

// FindLastHandler обрабатывает запрос (GET /news/last).
func (h *Handler) FindLastHandler(c *fiber.Ctx) error {
	out, err := h.findLastUC.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Err(err))
	}

	post := MapPostToPostDTO(out)
	resp := FindLastResponse{
		Post: post,
	}

	return c.Status(fiber.StatusOK).JSON(api.Resp(resp))
}
