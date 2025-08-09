// Package handler содержит http Обработчики
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type ServiceRequest struct {
	RouteName string
	Path      string
}

func (h *Handler) handleServiceRequest(c *fiber.Ctx, r ServiceRequest) error {
	body, status, err := h.fetchProxyResponse(c, r)
	if err != nil {
		return c.Status(status).JSON(
			fiber.Map{
				"status":  "error",
				"message": err.Error(),
			},
		)
	}

	var raw json.RawMessage
	if err = json.Unmarshal(body, &raw); err != nil {
		// Если невалидный JSON — возвращаем как есть (например, plain text)
		return c.Status(status).Send(body)
	}

	if status == fiber.StatusOK {
		return c.Status(status).JSON(
			fiber.Map{
				"data": raw,
			},
		)
	}

	return c.Status(status).JSON(raw)
}

func (h *Handler) fetchProxyResponse(c *fiber.Ctx, r ServiceRequest) ([]byte, int, error) {
	url, err := h.buildProxyURL(r.RouteName, r.Path, c)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	req, err := h.createProxyRequest(c, url)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fiber.StatusBadGateway, fmt.Errorf("failed to reach service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	c.Set(fiber.HeaderContentType, resp.Header.Get(fiber.HeaderContentType))

	return body, resp.StatusCode, nil
}

func (h *Handler) buildProxyURL(routeName, path string, c *fiber.Ctx) (string, error) {
	route, ok := h.registry.GetRouteByName(routeName)
	if !ok {
		return "", fmt.Errorf("service route not found")
	}

	url := route.BaseURL + path
	if query := c.Context().QueryArgs().String(); len(query) > 0 {
		url += "?" + string(query)
	}

	return url, nil
}

func (h *Handler) createProxyRequest(c *fiber.Ctx, url string) (*http.Request, error) {
	req, err := http.NewRequest(c.Method(), url, bytes.NewReader(c.Body()))
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy request")
	}

	// Прокидываем request ID, если есть
	if requestID, ok := c.Locals("request_id").(string); ok && requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	// Копируем заголовки из Fiber запроса
	for key, values := range c.GetReqHeaders() {
		if strings.ToLower(key) == "host" {
			continue
		}
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	return req, nil
}
