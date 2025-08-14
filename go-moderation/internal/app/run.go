package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ee-crocush/go-news/go-moderation/internal/adapter"
	"github.com/ee-crocush/go-news/go-moderation/internal/infrastructure/config"
	"github.com/ee-crocush/go-news/go-moderation/internal/service"
	"github.com/ee-crocush/go-news/pkg/kafka"
	"github.com/ee-crocush/go-news/pkg/logger"
	"github.com/rs/zerolog"
)

func Run(cfg *config.Config) {
	log := logger.GetLogger()

	// Создаем Publisher
	pub := initPublisher(cfg, log)
	defer pub.Close()

	// Создаем Consumer
	consumer := initConsumer(cfg, log, pub)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gracefulShutdown(ctx, consumer, log)
}

func initPublisher(cfg *config.Config, log *zerolog.Logger) *kafka.Publisher {
	topic, err := cfg.GetTopic("comment_moderated")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get topic")
	}

	return kafka.NewPublisher(cfg.Kafka.Brokers, topic)
}

func initConsumer(cfg *config.Config, log *zerolog.Logger, pub *kafka.Publisher) *kafka.Consumer {
	moderationService := service.NewService(pub)

	topic, err := cfg.GetTopic("comment_created")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get topic")
	}

	moderationAdapter := adapter.NewModerationAdapter(moderationService)
	return kafka.NewConsumer(cfg.Kafka.Brokers, topic, cfg.Kafka.ConsumerGroup, moderationAdapter)
}

func gracefulShutdown(ctx context.Context, consumer *kafka.Consumer, log *zerolog.Logger) {
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
	fmt.Println("Shutting down moderation service...")
}
