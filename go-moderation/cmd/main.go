// Package main представляет собой точку входа для приложения.
package main

import (
	"fmt"
	"github.com/ee-crocush/go-news/go-moderation/internal/app"
	"github.com/ee-crocush/go-news/go-moderation/internal/infrastructure/config"
	"github.com/ee-crocush/go-news/pkg/logger"
)

func main() {
	configPath := "./configs/config.yaml"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	logger.InitLogger(cfg.App.Name)

	app.Run(cfg)
}
