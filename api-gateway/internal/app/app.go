// Package app запускает сервер.
package app

import (
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/config"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/registry"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib"
	"github.com/ee-crocush/go-news/pkg/server"
	commonFiber "github.com/ee-crocush/go-news/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
	"time"
)

// Run запускает HTTP сервер и инициализирует все необходимые компоненты.
func Run(cfg *config.Config) error {
	services := registry.NewRouteRegistry(cfg.Routes)
	timeout := time.Duration(cfg.App.ConnectTimeout) * time.Second
	handlers := httplib.NewHandlers(cfg, services, timeout)

	// Создаем Fiber сервер
	fiberServer := commonFiber.NewFiberServer(
		cfg, func(app *fiber.App) {
			httplib.SetupRoutes(app, handlers)
		},
	)

	serverManager := server.NewServerManager(fiberServer)

	return serverManager.StartAll()
}
