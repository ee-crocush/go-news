// Package main содержит точку входа в приложение.
package main

import (
	"fmt"

	"github.com/ee-crocush/go-news/go-comments/internal/app"
	"github.com/ee-crocush/go-news/go-comments/internal/infrastructure/config"
	configLoader "github.com/ee-crocush/go-news/pkg/config"
	"github.com/ee-crocush/go-news/pkg/logger"
)

func main() {
	configPath := configLoader.FindConfigFile()
	cfg, err := config.LoadConfig(configPath)

	if err != nil || cfg == nil {
		fmt.Println("failed to load config from all known paths:", err)
		return
	}

	logger.InitLogger(cfg.App.Name)

	if err = app.Run(cfg); err != nil {
		fmt.Println("service failed to start:", err)
	}
}
