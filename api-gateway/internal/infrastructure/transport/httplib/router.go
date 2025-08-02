// Package httplib управляет настройкой маршрутов HTTP.
package httplib

import (
	_ "github.com/ee-crocush/go-news/api-gateway/docs"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/service"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/handler/comments"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/handler/health"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/handler/news"
	fiberServer "github.com/ee-crocush/go-news/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"time"
)

// Handlers группирует все обработчики.
type Handlers struct {
	News     *news.Handler
	Comments *comments.Handler
	Health   *health.Handler
}

// NewHandlers создает все обработчики.
func NewHandlers(cfg fiberServer.Config, registry service.RegistryService, timeout time.Duration) *Handlers {
	return &Handlers{
		News:     news.NewHandler(registry, timeout),
		Comments: comments.NewHandler(),
		Health:   health.NewHandler(cfg, registry, timeout),
	}
}

// SetupRoutes регистрирует маршруты.
func SetupRoutes(app *fiber.App, handlers *Handlers) {
	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault)
	// Health checks для самого API Gateway
	app.Get("/health", handlers.Health.HealthCheckHandler)
	app.Get("/ready", handlers.Health.ReadinessHandler)

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
		newsGroup.Get("/", h.FindAll)
		newsGroup.Get("/last", h.FindLast)
		newsGroup.Get("/latest/:limit?", h.FindLatest)
		newsGroup.Get("/:id", h.FindByID)
	}
}

// setupCommentsRoutes настраивает маршруты для комментариев.
func setupCommentsRoutes(api fiber.Router, h *comments.Handler) {
	commentsGroup := api.Group("/comments")
	{
		// Получение комментариев к новости
		commentsGroup.Get("/news/:newsId", h.GetCommentsHandler)
		commentsGroup.Post("/", h.CreateCommentHandler)
	}
}
