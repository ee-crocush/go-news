package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
)

// AppConfig - конфигурация приложения.
type AppConfig struct {
	Name    string `yaml:"name" validate:"required"`
	Version string `yaml:"version" validate:"required"`
}

// LoggingConfig - конфигурация логирования.
type LoggingConfig struct {
	Level  string `yaml:"level" validate:"required"`
	Format string `yaml:"format" validate:"required"`
}

// KafkaConfig - конфигурация Kafka.
type KafkaConfig struct {
	Brokers       []string          `yaml:"brokers" validate:"required"`
	Topics        map[string]string `yaml:"topics" validate:"required"`
	ConsumerGroup string            `yaml:"consumer_group" validate:"required"`
}

// Config основная конфигурация.
type Config struct {
	App     AppConfig     `yaml:"app"`
	Logging LoggingConfig `yaml:"logging"`
	Kafka   KafkaConfig   `yaml:"kafka"`
}

func (c *Config) GetAppName() string {
	return c.App.Name
}

func (c *Config) GetVersion() string {
	return c.App.Version
}

func (c *Config) GetTopic(name string) (string, error) {
	if topic, ok := c.Kafka.Topics[name]; ok {
		return topic, nil
	}
	return "", fmt.Errorf("topic %s not found", name)
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
			if err = godotenv.Load("./go-moderation/.env"); err != nil {
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

	return &cfg, nil
}
