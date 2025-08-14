// Package config содержит конфигурацию приложения.
package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"time"
)

// AppConfig - конфигурация приложения.
type AppConfig struct {
	Name                string `yaml:"name" validate:"required"`
	Env                 string `yaml:"env" validate:"required"`
	Version             string `yaml:"version" validate:"required"`
	ReadTimeout         int    `yaml:"read_timeout" validate:"required"`
	WriteTimeout        int    `yaml:"write_timeout" validate:"required"`
	ConnectTimeout      int    `yaml:"connect_timeout" validate:"required"`
	EnableRequestID     bool   `yaml:"enable_request_id" validate:"required"`
	EnableLogging       bool   `yaml:"enable_logging" validate:"required"`
	EnableErrorHandling bool   `yaml:"enable_error_handling" validate:"required"`
	EnableCors          bool   `yaml:"enable_cors" validate:"required"`
}

// HTTPConfig - конфигурация HTTP сервера.
type HTTPConfig struct {
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"required"`
}

// LoggingConfig - конфигурация логирования.
type LoggingConfig struct {
	Level          string `yaml:"level" validate:"required"`
	Format         string `yaml:"format" validate:"required"`
	EnableHTTPLogs bool   `yaml:"enable_http_logs" validate:"required"`
}

// Route - конфигурация сервисов для роутинга и проверки состояния.
type Route struct {
	Name       string `yaml:"name" validate:"required"`
	BaseURL    string `yaml:"base_url" validate:"required"`
	HealthPath string `yaml:"health_path" validate:"required"`
}

// Config основная конфигурация.
type Config struct {
	App     AppConfig     `yaml:"app"`
	HTTP    HTTPConfig    `yaml:"http"`
	Logging LoggingConfig `yaml:"logging"`
	Routes  []Route       `yaml:"routes"`
}

func (c *Config) IsDevelopment() bool {
	return c.App.Env == "dev"
}

func (c *Config) GetAppName() string {
	return c.App.Name
}

func (c *Config) GetVersion() string {
	return c.App.Version
}

func (c *Config) GetHost() string {
	return c.HTTP.Host
}

func (c *Config) GetPort() int {
	return c.HTTP.Port
}

func (c *Config) GetReadTimeout() time.Duration {
	return time.Duration(c.App.ReadTimeout) * time.Second
}

func (c *Config) GetWriteTimeout() time.Duration {
	return time.Duration(c.App.WriteTimeout) * time.Second
}

func (c *Config) EnableRequestID() bool {
	return c.App.EnableRequestID
}

func (c *Config) EnableLogging() bool {
	return c.App.EnableLogging
}

func (c *Config) EnableErrorHandling() bool {
	return c.App.EnableErrorHandling
}

func (c *Config) EnableCors() bool {
	return c.App.EnableCors
}

// Validate валидация конфига.
func (c *Config) Validate() error {
	validate := validator.New()

	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("Config.Validate: %w", err)
	}

	return nil
}

// LoadConfig загружает конфиг из файла.
func LoadConfig(configPath string) (*Config, error) {
	appEnv := os.Getenv("APP_ENV")

	if appEnv != "prod" && appEnv != "production" {
		if err := godotenv.Load(); err != nil {
			if err = godotenv.Load("./api-gateway/.env"); err != nil {
				return nil, fmt.Errorf("error loading .env file: %w", err)
			}
		}
	}

	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	// Подставляем переменные окружения
	expanded := os.ExpandEnv(string(raw))

	// Парсим YAML
	var cfg Config
	if err = yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("parse config yaml: %w", err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	cfg.overrideRoutesForDev()

	return &cfg, nil
}

// overrideRoutesForDev переписывает пути для локальной разработки.
// Если переменная окружения APP_ENV равна "dev", то для всех роутов
// BaseURL будет изменен на http://localhost:<port>.
// Например, если BaseURL был "http://example.com:8080",
// то он станет "http://localhost:8080".
func (c *Config) overrideRoutesForDev() {
	if c.App.Env != "dev" {
		return
	}

	for i, route := range c.Routes {
		u, err := url.Parse(route.BaseURL)
		if err != nil {
			continue
		}

		port := u.Port()
		if port == "" {
			port = "80"
		}

		c.Routes[i].BaseURL = fmt.Sprintf("http://localhost:%s", port)
	}
}
