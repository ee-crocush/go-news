package handler

import (
	"GoNews/pkg/api"
	"github.com/gofiber/fiber/v2"
)

// FindAllResponse представляет ответ на запрос получения всех постов.
type FindAllResponse struct {
	Posts []PostItem `json:"posts"`
}

// FindAllHandler обрабатывает запрос (GET /news).
func (h *Handler) FindAllHandler(c *fiber.Ctx) error {
	out, err := h.findAllUC.Execute(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Err(err))
	}

	posts := make([]PostItem, 0, len(out))
	for _, post := range out {
		posts = append(
			posts, PostItem{
				ID:      post.ID,
				Title:   post.Title,
				Content: post.Content,
				Link:    post.Link,
				PubTime: post.PubTime,
			},
		)
	}

	resp := FindAllResponse{
		Posts: posts,
	}

	return c.Status(fiber.StatusOK).JSON(api.Resp(resp))
}
