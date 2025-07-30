package health

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

// Handler представляет HTTP-handler для проверки жизнеспособности системы.
type Handler struct {
	startTime time.Time
}

// NewHandler создает новый экземпляр HTTP-handler для health checks.
func NewHandler() *Handler {
	return &Handler{
		startTime: time.Now(),
	}
}

// HealthCheckHandler обработчик для проверки жизнеспособности API Gateway
func (h *Handler) HealthCheckHandler(c *fiber.Ctx) error {
	uptime := time.Since(h.startTime)

	return c.JSON(
		fiber.Map{
			"status":    "OK",
			"service":   "api-gateway",
			"message":   "API Gateway is running",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"uptime":    uptime.String(),
			"version":   "1.0.0", // можно вынести в конфиг
		},
	)
}

// ReadinessHandler проверка готовности системы (можно добавить проверки зависимостей)
func (h *Handler) ReadinessHandler(c *fiber.Ctx) error {
	// TODO: Добавить проверки доступности микросервисов
	// isNewsServiceReady := h.checkNewsService()
	// isCommentsServiceReady := h.checkCommentsService()

	return c.JSON(
		fiber.Map{
			"status":  "OK",
			"service": "api-gateway",
			"message": "API Gateway is ready",
			"checks": fiber.Map{
				"news_service":     "not_implemented", // TODO: реальная проверка
				"comments_service": "not_implemented", // TODO: реальная проверка
			},
		},
	)
}

// LivenessHandler проверка живости системы
func (h *Handler) LivenessHandler(c *fiber.Ctx) error {
	return c.JSON(
		fiber.Map{
			"status":  "OK",
			"service": "api-gateway",
			"message": "API Gateway is alive",
		},
	)
}
