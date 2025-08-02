// Package main содержит точку входа в приложение.
package main

import (
	"fmt"
	"github.com/ee-crocush/go-news/api-gateway/internal/app"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/config"
	"github.com/ee-crocush/go-news/pkg/logger"
)

func main() {
	paths := []string{
		"./configs/dev.yaml",
		"./api-gateway/configs/dev.yaml",
	}

	var (
		cfg *config.Config
		err error
	)

	for _, path := range paths {
		cfg, err = config.LoadConfig(path)
		if err == nil {
			break
		}
	}

	if err != nil || cfg == nil {
		fmt.Println("failed to load config from all known paths:", err)
		return
	}

	logger.InitLogger(cfg.App.Name)
	if err = app.Run(cfg); err != nil {
		fmt.Println("service failed to start:", err)
		return
	}
}
