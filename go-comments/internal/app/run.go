package app

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/config"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/kafka"
	repo "github.com/ee-crocush/go-news/go-comments/internal/infrastructure/repo/postgres"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/transport/httplib"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/transport/httplib/handler"
	uc "github.com/ee-crocush/go-news/go-comments/internal/usecase/comment"
	"github.com/ee-crocush/go-news/pkg/logger"
	"github.com/ee-crocush/go-news/pkg/server"
	commonFiber "github.com/ee-crocush/go-news/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
)

// Run запускает HTTP сервер и инициализирует все необходимые компоненты.
func Run(cfg config.Config) error {
	repository, err := connectDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connectDB: %w", err)
	}

	consumer, err := initConsumer(cfg, repository)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.GetLogger()
	// Запуск Kafka consumer в отдельной горутине. Не будет возвращать ошибку и ломать весь сервис, т.к.
	// в будущем планируем добавить consumer, который будет проверять, какие комментарии еще ожидают модерации
	go func() {
		if err = consumer.Start(ctx); err != nil {
			log.Error().Err(err).Msg("consumer failed to start")
		}
	}()

	commentHandler, err := initHandler(cfg, repository)
	if err != nil {
		return fmt.Errorf("failed to init handler: %w", err)
	}

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

func initHandler(cfg config.Config, repository *repo.CommentRepository) (*handler.Handler, error) {
	topic, err := cfg.GetTopic("comment_created")
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	commentPublisher := kafka.NewPublisher(cfg.Kafka.Brokers, topic)
	commentCreateUC := uc.NewCreateUseCase(repository, commentPublisher)
	commentFindAllUC := uc.NewFindAllByNewsIDUseCase(repository)

	return handler.NewHandler(commentCreateUC, commentFindAllUC), nil
}

func initConsumer(cfg config.Config, repository *repo.CommentRepository) (*kafka.Consumer, error) {
	topic, err := cfg.GetTopic("comment_moderated")
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	updateStatusUC := uc.NewChangeStatusUseCase(repository)
	consumer := kafka.NewConsumer(cfg.Kafka.Brokers, topic, cfg.Kafka.ConsumerGroup, updateStatusUC)

	return consumer, nil
}
