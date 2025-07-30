// Package httplib управляет настройкой маршрутов HTTP.
package httplib

import (
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/handler/comments"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/handler/health"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/handler/news"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Handlers группирует все обработчики
type Handlers struct {
	News     *news.Handler
	Comments *comments.Handler
	Health   *health.Handler
}

// NewHandlers создает все обработчики
func NewHandlers() *Handlers {
	return &Handlers{
		News:     news.NewHandler(),
		Comments: comments.NewHandler(),
		Health:   health.NewHandler(),
	}
}

// SetupRoutes регистрирует маршруты.
func SetupRoutes(app *fiber.App, handlers *Handlers) {
	app.Use(recover.New())

	// Health checks для самого API Gateway
	app.Get("/health", handlers.Health.HealthCheckHandler)
	app.Get("/ready", handlers.Health.ReadinessHandler)
	app.Get("/live", handlers.Health.LivenessHandler)

	// Группа API маршрутов
	api := app.Group("/api")
	setupNewsRoutes(api, handlers.News)
	setupCommentsRoutes(api, handlers.Comments)

	app.Use(
		func(c *fiber.Ctx) error {
			return c.Status(404).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Route not found",
					"path":    c.Path(),
					"method":  c.Method(),
				},
			)
		},
	)
}

// setupNewsRoutes настраивает маршруты для новостей.
func setupNewsRoutes(api fiber.Router, h *news.Handler) {
	newsGroup := api.Group("/news")
	{
		newsGroup.Get("/", h.FindAllHandler)                  // GET /api/news
		newsGroup.Get("/last", h.FindLastHandler)             // GET /api/news/last
		newsGroup.Get("/latest/:limit?", h.FindLatestHandler) // GET /api/news/latest/10
		newsGroup.Get("/:id", h.FindByIDHandler)              // GET /api/news/123
	}
}

// setupCommentsRoutes настраивает маршруты для комментариев.
func setupCommentsRoutes(api fiber.Router, h *comments.Handler) {
	commentsGroup := api.Group("/comments")
	{
		// Получение комментариев к новости
		commentsGroup.Get("/news/:newsId", h.GetCommentsHandler) // GET /api/comments/news/123
		commentsGroup.Post("/", h.CreateCommentHandler)          // POST /api/comments
	}
}
