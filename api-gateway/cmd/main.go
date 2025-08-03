// Package main содержит точку входа в приложение.
package main

import (
	"fmt"
	"github.com/ee-crocush/go-news/api-gateway/internal/app"
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/config"
	"github.com/ee-crocush/go-news/pkg/logger"
)

// @title GoNews API Gateway
// @version 1.0
// @description API-шлюз для микросервисов агрегатора новостей и комментариев.
// @termsOfService https://go-news.example.com/terms/

// @contact.name Вымышленная команда поддержки GoNews
// @contact.url https://go-news.example.com/support
// @contact.email support@go-news.example.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @schemes http
func main() {
	paths := []string{
		"./configs/config.yaml",
		"./api-gateway/configs/config.yaml",
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
