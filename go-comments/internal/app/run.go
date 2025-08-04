package app

import (
	"fmt"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/config"
	repo "github.com/ee-crocush/go-news/go-comments/internal/infrastructure/repo/postgres"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/transport/httplib"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/transport/httplib/handler"
	uc "github.com/ee-crocush/go-news/go-comments/internal/usecase/comment"
	"github.com/ee-crocush/go-news/pkg/server"
	commonFiber "github.com/ee-crocush/go-news/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
)

// Run запускает HTTP сервер и инициализирует все необходимые компоненты.
func Run(cfg config.Config) error {
	repository, err := connectDB(cfg)
	if err != nil {
		return fmt.Errorf("connectDB: %w", err)
	}

	commentHandler := initHandler(repository)
	// Создаем Fiber сервер
	fiberServer := commonFiber.NewFiberServer(
		&cfg, func(app *fiber.App) {
			httplib.SetupRoutes(app, commentHandler)
		},
	)

	// Запускаем сервер
	serverManager := server.NewServerManager(fiberServer)
	return serverManager.StartAll()
}

func connectDB(cfg config.Config) (*repo.CommentRepository, error) {
	pgxPool, err := repo.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	repository := repo.NewCommentRepository(pgxPool)

	fmt.Printf(
		"PostgreSQL connected successfully! host=%s, port=%d, database=%s\n", cfg.DB.Host, cfg.DB.Port,
		cfg.DB.Name,
	)

	return repository, nil
}

func initHandler(repository *repo.CommentRepository) *handler.Handler {
	commentCreateUC := uc.NewCreateUseCase(repository)
	commentFindAllUC := uc.NewFindAllByNewsIDUseCase(repository)

	return handler.NewHandler(commentCreateUC, commentFindAllUC)
}
