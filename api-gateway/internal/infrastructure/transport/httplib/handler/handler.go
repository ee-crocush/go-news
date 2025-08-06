// Package handler содержит http Обработчики
package handler

import (
	"bytes"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/service"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"strings"
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
	// Получаем эндпоинт
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
	if query := c.Context().QueryArgs().String(); len(query) > 0 {
		url += "?" + string(query)
	}
	// Создаем запрос
	req, err := http.NewRequest(c.Method(), url, bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Failed to create proxy request",
			},
		)
	}

	requestID, ok := c.Locals("request_id").(string)
	if ok && requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	// Прокидываем заголовки из запроса Fiber → http.Request
	for key, values := range c.GetReqHeaders() {
		if strings.ToLower(key) == "host" {
			continue // не копируем Host
		}
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	// Явно прокидываем X-Request-ID (на всякий случай, если его не было выше)
	//if requestID := c.Get("X-Request-ID"); requestID != "" {
	//	req.Header.Set("X-Request-ID", requestID)
	//}

	// Отправляем запрос
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Failed to reach service",
			},
		)
	}
	defer resp.Body.Close()

	// Копируем статус и тело ответа
	c.Status(resp.StatusCode)

	// Копируем заголовки из ответа
	for key, values := range resp.Header {
		for _, v := range values {
			c.Set(key, v)
		}
	}

	// Копируем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Failed to read response",
			},
		)
	}

	return c.Send(body)
}
