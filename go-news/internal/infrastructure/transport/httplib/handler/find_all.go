package handler

import (
	uc "github.com/ee-crocush/go-news/go-news/internal/usecase/post"
	"github.com/ee-crocush/go-news/pkg/api"
	"github.com/gofiber/fiber/v2"
)

// FindAllResponse представляет ответ на запрос получения всех постов.
type FindAllResponse struct {
	Posts []PostItem `json:"posts"`
}

// FindAllHandler обрабатывает запрос (GET /news).
func (h *Handler) FindAllHandler(c *fiber.Ctx) error {
	search := c.Query("search")

	var out []uc.PostDTO
	var err error

	if search != "" {
		in := uc.FindByTitleSubstringInputDTO{Substring: search}
		out, err = h.findByTitleSubstrUC.Execute(c.Context(), in)
	} else {
		out, err = h.findAllUC.Execute(c.Context())
	}

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
