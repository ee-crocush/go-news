// Package health содержит обработчики жизнеспособности сервиса.
package health

import (
	"fmt"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/service"
	fiberServer "github.com/ee-crocush/go-news/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

// Handler представляет HTTP-handler для проверки жизнеспособности системы.
type Handler struct {
	startTime  time.Time
	config     fiberServer.Config
	registry   service.RegistryService
	httpClient *http.Client
}

// NewHandler создает новый экземпляр HTTP-handler для health checks.
func NewHandler(cfg fiberServer.Config, registry service.RegistryService, timeout time.Duration) *Handler {
	return &Handler{
		startTime:  time.Now(),
		config:     cfg,
		registry:   registry,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// HealthCheckHandler обработчик для проверки жизнеспособности API Gateway.
// @Summary Health check
// @Description Возвращает статус сервиса, аптайм и версию.
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *Handler) HealthCheckHandler(c *fiber.Ctx) error {
	uptime := time.Since(h.startTime)

	return c.JSON(
		fiber.Map{
			"status":  "OK",
			"service": h.config.GetAppName(),
			"version": h.config.GetVersion(),
			"uptime":  uptime.String(),
		},
	)
}

// ReadinessHandler проверка готовности системы (можно добавить проверки зависимостей)
func (h *Handler) ReadinessHandler(c *fiber.Ctx) error {
	checks := make(map[string]string)

	for _, route := range h.registry.GetAllRoutes() {
		url := fmt.Sprintf("%s%s", route.BaseURL, route.HealthPath)

		resp, err := h.httpClient.Get(url)
		if err != nil || resp.StatusCode >= 400 {
			checks[route.Name] = "unhealthy"
			continue
		}

		checks[route.Name] = "healthy"
	}

	return c.JSON(
		fiber.Map{
			"status":  "OK",
			"service": h.config.GetAppName(),
			"message": "All services is ready",
			"checks":  checks,
		},
	)
}
