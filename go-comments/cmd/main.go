// Package main содержит точку входа в приложение.
package main

import (
	"fmt"
	"github.com/ee-crocush/go-news/go-comments/internal/app"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/config"
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

	if err = app.Run(cfg); err != nil {
		fmt.Println("service failed to start:", err)
	}
}
