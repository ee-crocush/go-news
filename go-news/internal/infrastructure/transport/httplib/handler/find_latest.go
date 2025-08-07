package handler

import (
	"fmt"
	uc "github.com/ee-crocush/go-news/go-news/internal/usecase/post"
	"github.com/ee-crocush/go-news/pkg/api"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// FindLatestRequest - входные данные из тела запроса для получения последних n новостей.
type FindLatestRequest struct {
	Limit int `json:"limit"`
}

// FindLatestResponse представляет ответ на запрос получения последних n новостей.
type FindLatestResponse struct {
	News []PostDTO `json:"news"`
}

// FindLatestHandler обрабатывает запрос (GET /news/latest).
func (h *Handler) FindLatestHandler(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "0") // дефолт 0 или другой
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(api.Err(fmt.Errorf("invalid limit parameter")))
	}

	in := uc.FindLatestInputDTO{Limit: limit}
	out, err := h.findLatestUC.Execute(c.Context(), in)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Err(err))
	}

	news := MapNewsToNewsDTO(out)

	resp := FindLatestResponse{
		News: news,
	}

	return c.Status(fiber.StatusOK).JSON(api.Resp(resp))
}
