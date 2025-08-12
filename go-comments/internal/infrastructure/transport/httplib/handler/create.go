package handler

import (
	uc "github.com/ee-crocush/go-news/go-comments/internal/usecase/comment"
	"github.com/ee-crocush/go-news/pkg/api"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// CreateRequest - входные данные из тела запроса для создания комментария.
type CreateRequest struct {
	NewsID   int32  `json:"news_id" validate:"required,gt=0"`
	ParentID *int64 `json:"parent_id,omitempty"`
	Username string `json:"username" validate:"required,min=6,max=50"`
	Content  string `json:"content" validate:"required,min=1"`
}

// CreateRequestResponse представляет выходной данные запроса.
type CreateRequestResponse struct {
	Message string `json:"message"`
}

// CreateHandler обрабатывает запрос на создание нового комментария (Post /comments).
func (h *Handler) CreateHandler(c *fiber.Ctx) error {
	var req CreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(api.ErrWithCode("invalid-body", "Invalid request body"))
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(api.ErrWithCode("validation-error", validationErrors.Error()))
	}

	dto := uc.CommentDTO{
		NewsID:   req.NewsID,
		ParentID: req.ParentID,
		Username: req.Username,
		Content:  req.Content,
	}

	if err := h.createUC.Execute(c.Context(), dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Err(err))
	}

	response := CreateRequestResponse{
		Message: "Comment created successfully",
	}

	return c.Status(fiber.StatusCreated).JSON(api.Resp(response))
}
