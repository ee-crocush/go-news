// Package fiber представляет модуль реализации сервера на Fiber.
package fiber

import (
	"GoNews/internal/infrastructure/config"
	"GoNews/internal/infrastructure/transport/httplib"
	http_handler "GoNews/internal/infrastructure/transport/httplib/handler"
	http "GoNews/pkg/server"
	"context"

	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

// FiberServer реализует интерфейс GracefulServer.
type FiberServer struct {
	app *fiber.App
	cfg config.Config
}

// NewFiberServer создает Fiber сервер.
func NewFiberServer(cfg config.Config, h *http_handler.Handler) *FiberServer {
	app := http.NewFiberApp(&fiberConfig{Config: cfg})

	httplib.SetupRoutes(app, h)

	return &FiberServer{app: app, cfg: cfg}
}

// Start запускает Fiber.
func (s *FiberServer) Start() error {
	address := fmt.Sprintf("%s:%d", s.cfg.HTTP.Host, s.cfg.HTTP.Port)
	fmt.Printf("Fiber server listening on %s\n", address)
	return s.app.Listen(address)
}

// Shutdown останавливает Fiber-приложение.
func (s *FiberServer) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down Fiber server...")

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

// fiberConfig — отдельная структура для передачи настроек в Fiber.
type fiberConfig struct {
	config.Config
}

func (f *fiberConfig) GetAppName() string { return f.Config.App.Name }
func (f *fiberConfig) GetReadTimeout() time.Duration {
	return time.Duration(f.Config.App.ReadTimeout) * time.Second
}
func (f *fiberConfig) GetWriteTimeout() time.Duration {
	return time.Duration(f.Config.App.WriteTimeout) * time.Second
}
func (f *fiberConfig) EnableRequestID() bool     { return f.Config.App.EnableRequestID }
func (f *fiberConfig) EnableLogging() bool       { return f.Config.App.EnableLogging }
func (f *fiberConfig) EnableErrorHandling() bool { return f.Config.App.EnableErrorHandling }
