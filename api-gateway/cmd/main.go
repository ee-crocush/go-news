// Package main содержит точку входа в приложение.
package main

import (
	"fmt"
	"github.com/ee-crocush/go-news/api-gateway/internal/app"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/config"
	"github.com/ee-crocush/go-news/pkg/logger"
)

func main() {
	configPath := "./api-gateway/configs/dev.yaml"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("failed to load config:", err)
		return
	}

	logger.InitLogger(cfg.App.Name)
	if err = app.Run(cfg); err != nil {
		fmt.Println("service failed to start:", err)
		//log.Fatal().Err(err).Msg("Service failed to start")
		return
	}
}
