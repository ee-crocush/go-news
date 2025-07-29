// Package main содержит точку входа в приложение.
package main

import (
	"GoNews/internal/app"
	"GoNews/internal/infrastructure/config"
	"GoNews/pkg/logger"
	"fmt"
)

func main() {
	configPath := "./configs/dev.yaml"
	rssConfigPath := "./configs/rss_config.json"
	cfg, err := config.LoadConfig(configPath, rssConfigPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	logger.InitLogger(cfg.App.Name)
	log := logger.GetLogger()

	if err = app.Run(cfg, log); err != nil {
		log.Fatal().Err(err).Msg("Service failed to start")
	}
}
