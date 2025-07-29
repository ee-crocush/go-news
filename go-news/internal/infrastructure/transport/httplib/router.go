// Package httplib управляет настройкой маршрутов HTTP.
package httplib

import (
	"GoNews/internal/infrastructure/transport/httplib/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes регистрирует маршруты для Fiber приложения.
func SetupRoutes(app *fiber.App, h *handler.Handler) {
	//app.Static("/", "./public")
	app.Static(
		"/", "./public", fiber.Static{
			Index:    "index.html",
			Browse:   false,
			MaxAge:   86400,
			Compress: true,
		},
	)
	app.Get("/health", h.HealthCheckHandler)
	app.Get("/news", h.FindAllHandler)
	app.Get("/news/last", h.FindLastHandler)
	app.Get("/news/latest/:limit?", h.FindLatestHandler)
	app.Get("/news/:id", h.FindByIDHandler)
}
