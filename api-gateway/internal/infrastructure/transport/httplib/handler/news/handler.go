// Package news содержит все обработчики HTTP запросов к сервису go-news
package news

import (
	"fmt"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/service"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"time"
)

// Handler представляет HTTP-handler для работы с новостями через API Gateway.
type Handler struct {
	registry   service.RegistryService
	httpClient *http.Client
}

// NewHandler создает новый экземпляр HTTP-handler.
func NewHandler(registry service.RegistryService, timeout time.Duration) *Handler {
	return &Handler{
		registry:   registry,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// proxyRequest проксирует запросы к сервису.
func (h *Handler) proxyRequest(c *fiber.Ctx, routeName, path string) error {
	route, ok := h.registry.GetRouteByName(routeName)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Service route not found",
			},
		)
	}

	url := route.BaseURL + path
	resp, err := h.httpClient.Get(url)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Failed to reach service",
			},
		)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.Status(resp.StatusCode).Send(body)
}

// FindAll получает все новости.
// @Summary Получить все новости
// @Description Возвращает список всех новостей.
// @Tags news
// @Produce json
// @Success 200 {array} PostResponse
// @Router /api/news [get]
func (h *Handler) FindAll(c *fiber.Ctx) error {
	return h.proxyRequest(c, "go-news", "/news")
}

// FindLast получает последнюю новость.
// @Summary Получить последнюю новость
// @Description Возвращает последнюю новость.
// @Tags news
// @Produce json
// @Success 200 {object} PostResponse
// @Router /api/news/last [get]
func (h *Handler) FindLast(c *fiber.Ctx) error {
	return h.proxyRequest(c, "go-news", "/news/last")
}

// FindLatest получает последние n новости.
// @Summary Получить последние n новостей
// @Description Возвращает последние n новостей.
// @Tags news
// @Param limit path int false "Количество последних новостей" default(10)
// @Produce json
// @Success 200 {object} PostResponse
// @Router /api/news/latest/{limit} [get]
func (h *Handler) FindLatest(c *fiber.Ctx) error {
	limit := c.Params("limit", "10")
	path := fmt.Sprintf("/news/latest/?=%s", limit)

	return h.proxyRequest(c, "go-news", path)
}

// FindByID получает новость по ID.
// @Summary Получить новость по ID
// @Description Возвращает новость по ID.
// @Tags news
// @Param id path string true "ID новости"
// @Produce json
// @Success 200 {object} PostResponse
// @Router /api/news/{id} [get]
func (h *Handler) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")
	path := fmt.Sprintf("/news/%s", id)

	return h.proxyRequest(c, "go-news", path)
}
