package middleware

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"io"
)

// LoggingMiddleware возвращает middleware для логирования HTTP-запросов в Fiber.
func LoggingMiddleware(log zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body string
		if c.Request().Header.ContentLength() > 0 {
			body = string(c.Body())
		}

		queryParams := c.Context().QueryArgs().String()
		ip := c.IP()

		// Получаем ответ
		var responseBody bytes.Buffer
		writer := io.MultiWriter(c.Response().BodyWriter(), &responseBody)
		c.Response().SetBodyStreamWriter(
			func(w *bufio.Writer) {
				mw := bufio.NewWriter(writer)
				defer mw.Flush()
				_, _ = io.Copy(mw, c.Response().BodyStream())
			},
		)

		ctx, done := logger.LogRequest(
			log, c.Context(), "http", c.Method(), c.Path(),
		)
		c.SetUserContext(ctx)

		err := c.Next()
		status := c.Response().StatusCode()
		evt := log.With().
			Str("request_id", ctx.Value("request_id").(string)).
			Int("status_code", status).
			Str("ip", ip).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status_code", status).
			Str("query_params", queryParams).
			Str("body", body).
			Str("response", truncate(responseBody.String(), 1000)).
			Logger()

		switch {
		case err != nil:
			evt.Error().Err(err).Msg("Request failed")
		case status >= 500:
			evt.Error().Msg("Request completed with server error")
		case status >= 400:
			evt.Warn().Msg("Request completed with client error")
		default:
			evt.Info().Msg("Request completed")
		}

		done()
		return err
	}
}

// RequestIDMiddleware генерирует и добавляет request_id.
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID", uuid.New().String())
		if existingID := c.Locals("request_id"); existingID != nil {
			requestID = existingID.(string)
		}
		c.Set("X-Request-ID", requestID)
		c.Locals("request_id", requestID)

		// Добавляем request_id в контекст Go
		ctx := context.WithValue(c.Context(), "request_id", requestID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

// ErrorHandlerMiddleware глобально обрабатывает ошибки.
func ErrorHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Locals("request_id").(string)

				logger.GetLogger().Error().
					Str("request_id", requestID).
					Str("stack_trace", fmt.Sprintf("%+v", err)).
					Msg("Recovered from panic")

				c.Status(fiber.StatusInternalServerError).JSON(
					fiber.Map{
						"status":  "error",
						"message": "Internal Server Error",
					},
				)
			}
		}()
		return c.Next()
	}
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "...(truncated)"
	}

	return s
}
