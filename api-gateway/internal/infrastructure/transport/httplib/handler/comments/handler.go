package comments

import "github.com/gofiber/fiber/v2"

// Handler представляет HTTP-handler для работы с комментариями через API Gateway.
type Handler struct {
	// В будущем здесь будет клиент для общения с comments-service
}

// NewHandler создает новый экземпляр HTTP-handler для комментариев.
func NewHandler() *Handler {
	return &Handler{}
}

// GetCommentsHandler заглушка для получения комментариев к новости
func (h *Handler) GetCommentsHandler(c *fiber.Ctx) error {
	newsID := c.Params("news_id")

	// Параметры запроса (пагинация, сортировка и т.д.)
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	sortBy := c.Query("sort", "created_at")

	// TODO: Проксировать запрос к comments-service
	return c.JSON(
		fiber.Map{
			"status":   "OK",
			"service":  "comments-service-proxy",
			"message":  "GetComments endpoint - will proxy to comments-service",
			"endpoint": "/comments/" + newsID,
			"params": fiber.Map{
				"newsId": newsID,
				"page":   page,
				"limit":  limit,
				"sort":   sortBy,
			},
			"data": []interface{}{},
		},
	)
}

// CreateCommentHandler заглушка для создания комментария
func (h *Handler) CreateCommentHandler(c *fiber.Ctx) error {
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
				"error":   err.Error(),
			},
		)
	}

	// TODO: Проксировать запрос к comments-service
	return c.Status(201).JSON(
		fiber.Map{
			"status":       "OK",
			"service":      "comments-service-proxy",
			"message":      "CreateComment endpoint - will proxy to comments-service",
			"endpoint":     "/comments",
			"request_body": requestBody,
			"data": fiber.Map{
				"id":         "stub-comment-id-123",
				"created_at": "2024-01-01T12:00:00Z",
			},
		},
	)
}
