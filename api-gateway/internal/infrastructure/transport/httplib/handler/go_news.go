package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/dto"
	"github.com/gofiber/fiber/v2"
)

// FindAllNews получает все новости.
// @Summary Получить все новости
// @Description Возвращает список всех новостей.
// @Tags news
// @Produce json
// @Success 200 {array} domain.Post
// @Router /api/news [get]
func (h *Handler) FindAllNews(c *fiber.Ctx) error {
	return h.handleServiceRequest(
		c, ServiceRequest{
			RouteName: NewsRouteName,
			Path:      "/news",
		},
	)
}

// FindLastNews получает последнюю новость.
// @Summary Получить последнюю новость
// @Description Возвращает последнюю новость.
// @Tags news
// @Produce json
// @Success 200 {object} domain.Post
// @Router /api/news/last [get]
func (h *Handler) FindLastNews(c *fiber.Ctx) error {
	return h.handleServiceRequest(
		c, ServiceRequest{
			RouteName: NewsRouteName,
			Path:      "/news/last",
		},
	)
}

// FindLatestNews получает последние n новости.
// @Summary Получить последние n новостей
// @Description Возвращает последние n новостей.
// @Tags news
// @Param limit path int false "Количество последних новостей" default(10)
// @Produce json
// @Success 200 {object} domain.Post
// @Router /api/news/latest/{limit} [get]
func (h *Handler) FindLatestNews(c *fiber.Ctx) error {
	limit := c.Params("limit", "10")
	path := fmt.Sprintf("/news/latest/?=%s", limit)

	return h.handleServiceRequest(
		c, ServiceRequest{
			RouteName: NewsRouteName,
			Path:      path,
		},
	)
}

// FindByIDNews получает новость по ID.
// @Summary Получить новость по ID
// @Description Возвращает новость по ID.
// @Tags news
// @Param id path string true "ID новости"
// @Produce json
// @Success 200 {object} domain.PostWithComments
// @Router /api/news/{id} [get]
func (h *Handler) FindByIDNews(c *fiber.Ctx) error {
	id := c.Params("id")

	// Запрос к сервису новостей
	newsService := ServiceRequest{RouteName: NewsRouteName, Path: fmt.Sprintf("/news/%s", id)}
	newsBody, status, err := h.fetchProxyResponse(c, newsService)
	if err != nil {
		return c.Status(status).JSON(
			fiber.Map{
				"status":  "error",
				"message": err.Error(),
			},
		)
	}
	if status >= 400 {
		return c.Send(newsBody)
	}

	var newsResp struct {
		Post dto.Post `json:"post"`
	}

	if err = json.Unmarshal(newsBody, &newsResp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Failed to parse news data",
			},
		)
	}

	// Запрос к сервису комментариев
	commentService := ServiceRequest{RouteName: CommentsRouteName, Path: fmt.Sprintf("/comments/news/%s", id)}
	commentBody, status, err := h.fetchProxyResponse(c, commentService)
	if err != nil {
		return c.Status(status).JSON(
			fiber.Map{
				"status":  "error",
				"message": err.Error(),
			},
		)
	}
	if status >= 400 {
		return c.Send(commentBody)
	}

	var commentsResp struct {
		Comments []dto.Comment `json:"comments"`
	}

	if err = json.Unmarshal(commentBody, &commentsResp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Failed to parse comment data",
			},
		)
	}

	response := dto.FindByIDResponse{
		Data: dto.PostWithComments{
			Post:     newsResp.Post,
			Comments: commentsResp.Comments,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
