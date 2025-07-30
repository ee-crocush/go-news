// Package handler содержит все обработчики HTTP запросов
package news

import "github.com/gofiber/fiber/v2"

// Handler представляет HTTP-handler для работы с новостями через API Gateway.
type Handler struct {
	// В будущем здесь будут клиенты для общения с микросервисами
}

// NewHandler создает новый экземпляр HTTP-handler.
func NewHandler() *Handler {
	return &Handler{}
}

// FindAllHandler заглушка для получения всех новостей
func (h *Handler) FindAllHandler(c *fiber.Ctx) error {
	// TODO: Проксировать запрос к news-service
	return c.JSON(
		fiber.Map{
			"status":   "OK",
			"service":  "news-service-proxy",
			"message":  "FindAll endpoint - will proxy to news-service",
			"endpoint": "/news",
			"data":     []interface{}{},
		},
	)
}

// FindLastHandler заглушка для получения последней новости
func (h *Handler) FindLastHandler(c *fiber.Ctx) error {
	// TODO: Проксировать запрос к news-service
	return c.JSON(
		fiber.Map{
			"status":   "OK",
			"service":  "news-service-proxy",
			"message":  "FindLast endpoint - will proxy to news-service",
			"endpoint": "/news/last",
			"data":     nil,
		},
	)
}

// FindLatestHandler заглушка для получения последних n новостей
func (h *Handler) FindLatestHandler(c *fiber.Ctx) error {
	limit := c.Params("limit", "10") // значение по умолчанию

	// TODO: Проксировать запрос к news-service
	return c.JSON(
		fiber.Map{
			"status":   "OK",
			"service":  "news-service-proxy",
			"message":  "FindLatest endpoint - will proxy to news-service",
			"endpoint": "/news/latest/" + limit,
			"limit":    limit,
			"data":     []interface{}{},
		},
	)
}

// FindByIDHandler заглушка для получения новости по ID
func (h *Handler) FindByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	// TODO: Проксировать запрос к news-service
	return c.JSON(
		fiber.Map{
			"status":   "OK",
			"service":  "news-service-proxy",
			"message":  "FindByID endpoint - will proxy to news-service",
			"endpoint": "/news/" + id,
			"id":       id,
			"data":     nil,
		},
	)
}
