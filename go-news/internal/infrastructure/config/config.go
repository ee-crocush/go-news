package config

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"net/url"
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

// MongoConfig конфигурация MongoDB.
type MongoConfig struct {
	Host           string        `yaml:"host" validate:"required"`
	Port           int           `yaml:"port" validate:"required"`
	User           string        `yaml:"user" validate:"required"`
	Password       string        `yaml:"password" validate:"required"`
	AuthSource     string        `yaml:"auth_source" validate:"required"`
	Database       string        `yaml:"database" validate:"required"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" validate:"required"`
}

// URI формирование строки подключения к БД.
func (c *MongoConfig) URI() *url.URL {
	hostPost := fmt.Sprintf("%s:%d", c.Host, c.Port)

	return &url.URL{
		Scheme: "mongodb",
		//User:   url.UserPassword(c.User, c.Password),
		Host: hostPost,
	}
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

// RSSConfig - конфигурация RSS источников.
type RSSConfig struct {
	RSS           []string `json:"rss" validate:"required,min=1,dive,url"`
	RequestPeriod int      `json:"request_period" validate:"required,min=1"`
}

func (r *RSSConfig) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return fmt.Errorf("RSSConfig.Validate: %w", err)
	}

	return nil
}

// GetRequestPeriodDuration возвращает период запросов как time.Duration в минутах.
func (r *RSSConfig) GetRequestPeriodDuration() time.Duration {
	return time.Duration(r.RequestPeriod) * time.Minute
}

// Config основная конфигурация.
type Config struct {
	App     AppConfig     `yaml:"app"`
	HTTP    HTTPConfig    `yaml:"http"`
	MongoDB MongoConfig   `yaml:"mongodb"`
	Logging LoggingConfig `yaml:"logging"`
	RSS     RSSConfig     `json:"-"`
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
func LoadConfig(configPath, rssConfigPath string) (Config, error) {
	err := godotenv.Load()
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

	rssConfig, err := LoadRSSConfig(rssConfigPath)
	if err != nil {
		return Config{}, fmt.Errorf("load rss config failed: %w", err)
	}

	cfg.RSS = rssConfig

	if err = cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// LoadRSSConfig згружает конфиг для парсинга RSS.
func LoadRSSConfig(rssConfgiPath string) (RSSConfig, error) {
	raw, err := os.ReadFile(rssConfgiPath)
	if err != nil {
		return RSSConfig{}, fmt.Errorf("read rss config failed: %w", err)
	}

	var rssConfig RSSConfig
	if err = json.Unmarshal(raw, &rssConfig); err != nil {
		return RSSConfig{}, fmt.Errorf("parse rss config failed: %w", err)
	}

	if err = rssConfig.Validate(); err != nil {
		return RSSConfig{}, fmt.Errorf("validate rss config failed: %w", err)
	}

	return rssConfig, nil
}
