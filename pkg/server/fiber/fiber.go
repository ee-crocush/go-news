// Package fiber представляет модуль реализации сервера на Fiber.
package fiber

import (
	"context"
	"github.com/ee-crocush/go-news/pkg/logger"
	mw "github.com/ee-crocush/go-news/pkg/middleware"

	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

// Config - интерфейс для передачи конфигурации Fiber.
type Config interface {
	GetAppName() string
	GetHost() string
	GetPort() int
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	EnableRequestID() bool
	EnableLogging() bool
	EnableErrorHandling() bool
}

// RouteSetup функция для настройки маршрутов.
type RouteSetup func(*fiber.App)

// FiberServer реализует интерфейс GracefulServer.
type FiberServer struct {
	app *fiber.App
	cfg Config
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

// NewFiberServer создает Fiber сервер.
func NewFiberServer(cfg Config, setupRoutes RouteSetup) *FiberServer {
	app := NewFiberApp(cfg)

	if setupRoutes != nil {
		setupRoutes(app)
	}

	return &FiberServer{app: app, cfg: cfg}
}

// Start запускает Fiber.
func (s *FiberServer) Start() error {
	address := fmt.Sprintf("%s:%d", s.cfg.GetHost(), s.cfg.GetPort())
	fmt.Printf("Fiber server (%s) listening on %s\n", s.cfg.GetAppName(), address)
	return s.app.Listen(address)
}

// Shutdown останавливает Fiber-приложение.
func (s *FiberServer) Shutdown(ctx context.Context) error {
	fmt.Printf("Shutting down Fiber server (%s)...\n", s.cfg.GetAppName())

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.app.Shutdown()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
