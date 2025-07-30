package config

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

// AppConfig - конфигурация приложения.
type AppConfig struct {
	Name                string `yaml:"name" validate:"required"`
	ReadTimeout         int    `yaml:"read_timeout" validate:"required"`
	WriteTimeout        int    `yaml:"write_timeout" validate:"required"`
	EnableRequestID     bool   `yaml:"enable_request_id" validate:"required"`
	EnableLogging       bool   `yaml:"enable_logging" validate:"required"`
	EnableErrorHandling bool   `yaml:"enable_error_handling" validate:"required"`
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

// Config основная конфигурация.
type Config struct {
	App     AppConfig     `yaml:"app"`
	HTTP    HTTPConfig    `yaml:"http"`
	Logging LoggingConfig `yaml:"logging"`
}

func (c *Config) GetAppName() string {
	return c.App.Name
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

// Validate валидация конфига.
func (c *Config) Validate() error {
	validate := validator.New()

	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("Config.Validate: %w", err)
	}

	return nil
}

// LoadConfig загружает конфиг из файла.
func LoadConfig(configPath string) (Config, error) {
	err := godotenv.Load("./api-gateway/.env")
	if err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	raw, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	// Подставляем переменные окружения
	expanded := os.ExpandEnv(string(raw))

	// Парсим YAML
	var cfg Config
	if err = yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config yaml: %w", err)
	}

	if err = cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}
