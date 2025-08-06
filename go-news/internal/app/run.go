package app

import (
	"context"
	"fmt"
	"github.com/ee-crocush/go-news/go-news/internal/infrastructure/config"
	repo "github.com/ee-crocush/go-news/go-news/internal/infrastructure/repo/mongo"
	"github.com/ee-crocush/go-news/go-news/internal/infrastructure/rss"
	"github.com/ee-crocush/go-news/go-news/internal/infrastructure/transport/httplib"
	"github.com/ee-crocush/go-news/go-news/internal/infrastructure/transport/httplib/handler"
	uc "github.com/ee-crocush/go-news/go-news/internal/usecase/post"
	"github.com/ee-crocush/go-news/pkg/logger"
	"github.com/ee-crocush/go-news/pkg/server"
	commonFiber "github.com/ee-crocush/go-news/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

// Run запускает HTTP сервер и инициализирует все необходимые компоненты.
func Run(cfg config.Config) error {
	client, db, err := connectDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log := logger.GetLogger()
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			fmt.Println("Failed to disconnect MongoDB: %w", err)
		} else {
			fmt.Println("MongoDB disconnected successfully")
		}
	}()

	postRepo := repo.NewPostRepository(db, cfg.MongoDB.ConnectTimeout)
	rssParser := rss.NewParser(cfg.RSS.GetRequestPeriodDuration())
	postStoreUC := uc.NewParseAndStoreUseCase(postRepo, rssParser)
	go startRSSBackgroundJob(cfg, postStoreUC, log)

	postHandler := initHandler(postRepo)
	// Создаем Fiber сервер
	fiberServer := commonFiber.NewFiberServer(
		&cfg, func(app *fiber.App) {
			httplib.SetupRoutes(app, postHandler)
		},
	)

	// Запускаем сервер
	serverManager := server.NewServerManager(fiberServer)
	return serverManager.StartAll()
}

func connectDB(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	client, db, err := repo.Init(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	fmt.Printf(
		"MongoDB connected successfully! host=%s, port=%d, database=%s\n", cfg.MongoDB.Host, cfg.MongoDB.Port,
		cfg.MongoDB.Database,
	)

	return client, db, nil
}

func initHandler(repos *repo.PostRepository) *handler.Handler {
	findByIDUC := uc.NewFindByIDUseCase(repos)
	findLastUC := uc.NewFindLastUseCase(repos)
	findLatestUC := uc.NewFindLatestUseCase(repos)
	findAllUC := uc.NewFindAllUseCase(repos)
	findByTitleSubstrUC := uc.NewFindByTitleSubstring(repos)

	return handler.NewHandler(findByIDUC, findLastUC, findLatestUC, findAllUC, findByTitleSubstrUC)
}

func startRSSBackgroundJob(cfg config.Config, ucp uc.ParseAndStoreUseCase, log *zerolog.Logger) {
	period := cfg.RSS.GetRequestPeriodDuration()
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		for _, url := range cfg.RSS.RSS {
			go func(u string) {
				in := uc.ParseAndStoreInputDTO{URL: u}
				if err := ucp.Execute(context.Background(), in); err != nil {
					log.Error().Err(err).Str("url", u).Msg("rss parse failed")
				}
			}(url)
		}
		<-ticker.C
	}
}
