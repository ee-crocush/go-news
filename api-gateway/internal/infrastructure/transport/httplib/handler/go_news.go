package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

const NewsRouteName = "go-news"

// PostResponse описывает структуру ответа на /news
type PostResponse struct {
	ID      int32  `json:"id" example:"1"`
	Title   string `json:"title" example:"Example title"`
	Content string `json:"content" example:"Example Long Content"`
	Link    string `json:"link" example:"https://example.com/news/1"`
	PubTime string `json:"pub_time" example:"2025-06-26 10:00:43"`
}

// FindAllNews получает все новости.
// @Summary Получить все новости
// @Description Возвращает список всех новостей.
// @Tags news
// @Produce json
// @Success 200 {array} PostResponse
// @Router /api/news [get]
func (h *Handler) FindAllNews(c *fiber.Ctx) error {
	return h.proxyRequest(c, NewsRouteName, "/news")
}

// FindLastNews получает последнюю новость.
// @Summary Получить последнюю новость
// @Description Возвращает последнюю новость.
// @Tags news
// @Produce json
// @Success 200 {object} PostResponse
// @Router /api/news/last [get]
func (h *Handler) FindLastNews(c *fiber.Ctx) error {
	return h.proxyRequest(c, NewsRouteName, "/news/last")
}

// FindLatestNews получает последние n новости.
// @Summary Получить последние n новостей
// @Description Возвращает последние n новостей.
// @Tags news
// @Param limit path int false "Количество последних новостей" default(10)
// @Produce json
// @Success 200 {object} PostResponse
// @Router /api/news/latest/{limit} [get]
func (h *Handler) FindLatestNews(c *fiber.Ctx) error {
	limit := c.Params("limit", "10")
	path := fmt.Sprintf("/news/latest/?=%s", limit)

	return h.proxyRequest(c, NewsRouteName, path)
}

// FindByIDNews получает новость по ID.
// @Summary Получить новость по ID
// @Description Возвращает новость по ID.
// @Tags news
// @Param id path string true "ID новости"
// @Produce json
// @Success 200 {object} PostResponse
// @Router /api/news/{id} [get]
func (h *Handler) FindByIDNews(c *fiber.Ctx) error {
	id := c.Params("id")
	path := fmt.Sprintf("/news/%s", id)

	return h.proxyRequest(c, NewsRouteName, path)
}
