// Package httplib управляет настройкой маршрутов HTTP.
package httplib

import (
	_ "github.com/ee-crocush/go-news/api-gateway/docs"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/service"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/handler"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib/handler/health"
	fiberServer "github.com/ee-crocush/go-news/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"time"
)

// Handlers группирует все обработчики.
// Так как эндпоинтов не так много и логика обработки одинаковая, хендлеры объявлены в одном пакете
type Handlers struct {
	Health       *health.Handler
	NewsComments *handler.Handler
}

// NewHandlers создает все обработчики.
func NewHandlers(cfg fiberServer.Config, registry service.RegistryService, timeout time.Duration) *Handlers {
	return &Handlers{
		NewsComments: handler.NewHandler(registry, timeout),
		Health:       health.NewHandler(cfg, registry, timeout),
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
	setupNewsRoutes(api, handlers.NewsComments)
	setupCommentsRoutes(api, handlers.NewsComments)

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
func setupNewsRoutes(api fiber.Router, h *handler.Handler) {
	newsGroup := api.Group("/news")
	{
		newsGroup.Get("/", h.FindAllNews)
		newsGroup.Get("/last", h.FindLastNews)
		newsGroup.Get("/latest/:limit?", h.FindLatestNews)
		newsGroup.Get("/:id", h.FindByIDNews)
	}
}

// setupCommentsRoutes настраивает маршруты для комментариев.
func setupCommentsRoutes(api fiber.Router, h *handler.Handler) {
	commentsGroup := api.Group("/comments")
	{
		// Получение комментариев к новости
		commentsGroup.Get("/", h.FindAllCommentsByNewsID)
		commentsGroup.Post("/", h.CreateComments)
	}
}
