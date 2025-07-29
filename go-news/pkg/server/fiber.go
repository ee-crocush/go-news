package server

import (
	"GoNews/pkg/logger"
	mw "GoNews/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"time"
)

// Config - интерфейс для передачи конфигурации Fiber.
type Config interface {
	GetAppName() string
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	EnableRequestID() bool
	EnableLogging() bool
	EnableErrorHandling() bool
}

// NewFiberApp создаёт настроенное Fiber-приложение.
func NewFiberApp(config Config) *fiber.App {
	app := fiber.New(
		fiber.Config{
			AppName:      config.GetAppName(),
			ReadTimeout:  config.GetReadTimeout(),
			WriteTimeout: config.GetWriteTimeout(),
		},
	)

	log := logger.GetLogger()

	if config.EnableRequestID() {
		app.Use(mw.RequestIDMiddleware())
	}

	if config.EnableLogging() {
		app.Use(mw.LoggingMiddleware(*log))
	}

	if config.EnableErrorHandling() {
		app.Use(mw.ErrorHandlerMiddleware())
	}

	return app
}
