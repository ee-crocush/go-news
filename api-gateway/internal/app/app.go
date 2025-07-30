// Package app запускает сервер.
package app

import (
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/config"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/transport/httplib"
	"github.com/ee-crocush/go-news/pkg/server"
	commonFiber "github.com/ee-crocush/go-news/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
)

// Run запускает HTTP сервер и инициализирует все необходимые компоненты.
func Run(cfg config.Config) error {
	handlers := httplib.NewHandlers()

	// Создаем Fiber сервер
	fiberServer := commonFiber.NewFiberServer(
		&cfg, func(app *fiber.App) {
			httplib.SetupRoutes(app, handlers)
		},
	)

	// Запускаем сервер
	serverManager := server.NewServerManager(fiberServer)
	return serverManager.StartAll()
}
