package handler

import (
	"GoNews/pkg/api"
	"github.com/gofiber/fiber/v2"
)

// FindLastResponse представляет ответ на запрос получения последней новости.
type FindLastResponse struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	PubTime string `json:"pub_time"`
}

// FindLastHandler обрабатывает запрос (GET /news/last).
func (h *Handler) FindLastHandler(c *fiber.Ctx) error {
	out, err := h.findLastUC.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Err(err))
	}

	resp := FindLastResponse{
		ID:      out.ID,
		Title:   out.Title,
		Content: out.Content,
		Link:    out.Link,
		PubTime: out.PubTime,
	}

	return c.Status(fiber.StatusOK).JSON(api.Resp(resp))
}
