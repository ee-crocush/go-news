// Package main содержит точку входа в приложение.
package main

import (
	"fmt"

	"github.com/ee-crocush/go-news/go-news/internal/app"
	"github.com/ee-crocush/go-news/go-news/internal/infrastructure/config"
	configLoader "github.com/ee-crocush/go-news/pkg/config"
	"github.com/ee-crocush/go-news/pkg/logger"
)

func main() {
	configPath := configLoader.FindConfigFile()

	rssPathes := []string{
		"./configs/rss_config.json",
		"./go-news/configs/rss_config.json",
		"/app/configs/rss_config.json",
	}
	rssConfigPath := configLoader.FindConfigFile(rssPathes...)
	cfg, err := config.LoadConfig(configPath, rssConfigPath)
	if err != nil || cfg == nil {
		fmt.Println("failed to load config from all known paths:", err)
		return
	}

	logger.InitLogger(cfg.App.Name)

	if err = app.Run(cfg); err != nil {
		fmt.Println("service failed to start:", err)
	}
}
