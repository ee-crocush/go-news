package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) HealthCheckHandler(c *fiber.Ctx) error {
	err := c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"status":  "ok",
			"message": "Service is healthy",
		},
	)

	if err != nil {
		return fmt.Errorf("failed to send JSON response: %w", err)
	}

	return nil
}
