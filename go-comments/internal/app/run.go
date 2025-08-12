// Package app выполняет основную инициализацию сервиса.
package app

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/config"
	repo "github.com/ee-crocush/go-news/go-comments/internal/infrastructure/repo/postgres"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/transport/httplib"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/transport/httplib/handler"
	uc "github.com/ee-crocush/go-news/go-comments/internal/usecase/comment"
	"github.com/ee-crocush/go-news/pkg/kafka"
	"github.com/ee-crocush/go-news/pkg/logger"
	"github.com/ee-crocush/go-news/pkg/server"
	commonFiber "github.com/ee-crocush/go-news/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

// Run запускает HTTP сервер и инициализирует все необходимые компоненты.
func Run(cfg config.Config) error {
	repository, err := connectDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connectDB: %w", err)
	}

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

	log := logger.GetLogger()

	consumer, err := initConsumer(cfg, repository)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gracefulShutdownConsumer(ctx, consumer, log)

	// Запускаем сервер
	serverManager := server.NewServerManager(fiberServer)
	return serverManager.StartAll()
}

// connectDB выполняет подключение к БД.
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

// initHandler создает хендлеры.
func initHandler(cfg config.Config, repository *repo.CommentRepository) (*handler.Handler, error) {
	// Создаем топик, куда будем отправлять события о создании нового комментария, подлежащего модерации
	topic, err := cfg.GetTopic("comment_created")
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	commentPublisher := kafka.NewPublisher(cfg.Kafka.Brokers, topic)
	commentCreateUC := uc.NewCreateUseCase(repository, commentPublisher)
	commentFindAllUC := uc.NewFindAllByNewsIDUseCase(repository)

	return handler.NewHandler(commentCreateUC, commentFindAllUC), nil
}

// initConsumer создает consumer кафки для получения результатов модерации.
func initConsumer(cfg config.Config, repository *repo.CommentRepository) (*kafka.Consumer, error) {
	topic, err := cfg.GetTopic("comment_moderated")
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	updateStatusUC := uc.NewChangeStatusUseCase(repository)
	consumer := kafka.NewConsumer(cfg.Kafka.Brokers, topic, cfg.Kafka.ConsumerGroup, updateStatusUC)

	return consumer, nil
}

// gracefulShutdownConsumer реализация gracefulShutdown для kafka.Consumer
func gracefulShutdownConsumer(ctx context.Context, consumer *kafka.Consumer, log *zerolog.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatal().Err(err).Msg("Kafka consumer error")
		}
	}()

	<-sigChan
	fmt.Println("Shutting kafka consumer...")
	consumer.Close()
}
