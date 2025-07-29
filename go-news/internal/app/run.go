package app

import (
	"GoNews/internal/infrastructure/config"
	repo "GoNews/internal/infrastructure/repo/mongo"
	"GoNews/internal/infrastructure/rss"
	"GoNews/internal/infrastructure/server"
	"GoNews/internal/infrastructure/transport/httplib/handler"
	uc "GoNews/internal/usecase/post"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

// Run запускает HTTP сервер и инициализирует все необходимые компоненты.
func Run(cfg config.Config, log *zerolog.Logger) error {
	client, db, err := connectDB(log, cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Error().Err(err).Msg("Failed to disconnect MongoDB")
		} else {
			log.Info().Msg("MongoDB disconnected successfully")
		}
	}()

	postRepo := repo.NewPostRepository(db, cfg.MongoDB.ConnectTimeout)
	rssParser := rss.NewParser(cfg.RSS.GetRequestPeriodDuration())
	postStoreUC := uc.NewParseAndStoreUseCase(postRepo, rssParser)
	go startRSSBackgroundJob(cfg, postStoreUC, log)

	postHandler := initHandler(postRepo)

	servers, err := server.CreateServers(cfg, postHandler)
	if err != nil {
		return fmt.Errorf("failed to create servers: %w", err)
	}

	return server.StartAll(servers...)
}

func connectDB(log *zerolog.Logger, cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	client, db, err := repo.Init(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	log.Info().
		Str("host", cfg.MongoDB.Host).
		Int("port", cfg.MongoDB.Port).
		Str("database", cfg.MongoDB.Database).
		Msg("MongoDB connected successfully!")

	return client, db, nil
}

func initHandler(repo *repo.PostRepository) *handler.Handler {
	findByIDUC := uc.NewFindByIDUseCase(repo)
	findLastUC := uc.NewFindLastUseCase(repo)
	findLatestUC := uc.NewFindLatestUseCase(repo)
	findAllUC := uc.NewFindAllUseCase(repo)

	return handler.NewHandler(findByIDUC, findLastUC, findLatestUC, findAllUC)
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
