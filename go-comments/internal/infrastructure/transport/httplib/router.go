// Package httplib управляет настройкой маршрутов HTTP.
package httplib

import (
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/transport/httplib/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes регистрирует маршруты для Fiber приложения.
func SetupRoutes(app *fiber.App, h *handler.Handler) {
	app.Get("/health", h.HealthCheckHandler)

	commentsGroup := app.Group("/comments")
	{
		commentsGroup.Get("/news/:id", h.FindAllByNewsIDHandler)
		commentsGroup.Post("/", h.CreateHandler)
	}
}
