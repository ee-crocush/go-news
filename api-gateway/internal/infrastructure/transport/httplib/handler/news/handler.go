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

// FindAll получает все новости
func (h *Handler) FindAll(c *fiber.Ctx) error {
	return h.proxyRequest(c, "go-news", "/news")
}

// FindLast получает последнюю новость.
func (h *Handler) FindLast(c *fiber.Ctx) error {
	return h.proxyRequest(c, "go-news", "/news/last")
}

// FindLatest получает последние n новости.
func (h *Handler) FindLatest(c *fiber.Ctx) error {
	limit := c.Params("limit", "10")
	path := fmt.Sprintf("/news/latest/?=%s", limit)

	return h.proxyRequest(c, "go-news", path)
}

// FindByID получает новость по ID.
func (h *Handler) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")
	path := fmt.Sprintf("/news/%s", id)

	return h.proxyRequest(c, "go-news", path)
}
