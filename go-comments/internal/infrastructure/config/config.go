package config

import (
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
	Version             string `yaml:"version" validate:"required"`
	ReadTimeout         int    `yaml:"read_timeout" validate:"required"`
	WriteTimeout        int    `yaml:"write_timeout" validate:"required"`
	EnableRequestID     bool   `yaml:"enable_request_id" validate:"required"`
	EnableLogging       bool   `yaml:"enable_logging" validate:"required"`
	EnableErrorHandling bool   `yaml:"enable_error_handling" validate:"required"`
	EnableCors          bool   `yaml:"enable_cors"`
}

// DBConfig конфигурация базы данных.
type DBConfig struct {
	Host                string        `yaml:"host" validate:"required"`
	Port                int           `yaml:"port" validate:"required"`
	User                string        `yaml:"user" validate:"required"`
	Password            string        `yaml:"password" validate:"required"`
	Name                string        `yaml:"name" validate:"required"`
	Migrations          string        `yaml:"migrations" validate:"required"`
	SSLMode             string        `yaml:"sslmode" validate:"required"`
	PoolMaxConns        string        `yaml:"pool_max_conns" validate:"required"`
	PoolMinConns        string        `yaml:"pool_min_conns" validate:"required"`
	PoolMaxConnLifetime string        `yaml:"pool_max_conn_lifetime" validate:"required"`
	PoolMaxConnIdletime string        `yaml:"pool_max_conn_idle_time" validate:"required"`
	ConnectTimeout      time.Duration `yaml:"connect_timeout" validate:"required"`
}

// DSN формирование строки подключения к БД.
func (c *DBConfig) DSN() *url.URL {
	hostPost := fmt.Sprintf("%s:%d", c.Host, c.Port)

	return &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   hostPost,
		Path:   c.Name,
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

// Config основная конфигурация.
type Config struct {
	App     AppConfig     `yaml:"app"`
	HTTP    HTTPConfig    `yaml:"http"`
	DB      DBConfig      `yaml:"database"`
	Logging LoggingConfig `yaml:"logging"`
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
func LoadConfig(configPath, rssConfigPath string) (Config, error) {
	if err := godotenv.Load(); err != nil {
		if err = godotenv.Load("./go-news/.env"); err != nil {
			return Config{}, fmt.Errorf("error loading .env file: %w", err)
		}
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
