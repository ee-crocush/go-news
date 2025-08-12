package handler

import (
	uc "github.com/ee-crocush/go-news/go-news/internal/usecase/post"
	"github.com/ee-crocush/go-news/pkg/api"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// FindAllResponse представляет ответ на запрос получения всех постов.
type FindAllResponse struct {
	News  []PostDTO `json:"news"`
	Total int32     `json:"total"`
}

// FindAllHandler обрабатывает запрос (GET /news).
func (h *Handler) FindAllHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 10 {
		limit = 10
	}

	in := uc.FindAllInputDTO{
		Search: search,
		Limit:  limit,
		Page:   page,
	}
	out, total, err := h.findAllUC.Execute(c.Context(), in)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Err(err))
	}

	news := MapNewsToNewsDTO(out)

	resp := FindAllResponse{
		News:  news,
		Total: total,
	}

	return c.Status(fiber.StatusOK).JSON(api.Resp(resp))
}
