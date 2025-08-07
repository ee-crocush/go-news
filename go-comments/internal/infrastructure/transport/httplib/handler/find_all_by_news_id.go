package handler

import (
	uc "github.com/ee-crocush/go-news/go-comments/internal/usecase/comment"
	"github.com/ee-crocush/go-news/pkg/api"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// FindAllByNewsIDRequest - входные данные из тела запроса для получения всех комментариев конкретной новости.
type FindAllByNewsIDRequest struct {
	NewsID int32 `json:"news_id" validate:"required,gt=0"`
}

// FindAllByNewsIDResponse представляет ответ на запрос получения всех комментариев конкретной новости.
// Создавать еще отдельно структуру нет смысла, т.к. придется снова рекурсивно проходится по массиву.
type FindAllByNewsIDResponse struct {
	Comments []uc.CommentDTO `json:"comments"`
}

// FindAllByNewsIDHandler обрабатывает запрос на получение всех комментариев конкретного поста (GET /comments).
func (h *Handler) FindAllByNewsIDHandler(c *fiber.Ctx) error {
	var req FindAllByNewsIDRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(api.ErrWithCode("invalid-body", "Invalid request body"))
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(api.ErrWithCode("validation-error", validationErrors.Error()))
	}
	in := uc.AllByNewsIDDTO{NewsID: req.NewsID}
	out, err := h.findAllByNewsUC.Execute(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Err(err))
	}

	response := FindAllByNewsIDResponse{
		Comments: out,
	}

	return c.Status(fiber.StatusOK).JSON(api.Resp(response))
}
