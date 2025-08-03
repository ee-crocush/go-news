// Package main содержит точку входа в приложение.
package main

import (
	"fmt"
	"github.com/ee-crocush/go-news/go-news/internal/app"
	"github.com/ee-crocush/go-news/go-news/internal/infrastructure/config"
	"github.com/ee-crocush/go-news/pkg/logger"
)

func main() {
	configPath := "./configs/config.yaml"
	rssConfigPath := "./configs/rss_config.json"
	cfg, err := config.LoadConfig(configPath, rssConfigPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	logger.InitLogger(cfg.App.Name)

	if err = app.Run(cfg); err != nil {
		fmt.Println("service failed to start:", err)
	}
}
