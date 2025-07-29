// Package logger представляет модуль логирования.
package logger

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"os"
	"sync"
	"time"
)

var log zerolog.Logger
var once sync.Once

// InitLogger инициализирует глобальный логгер.
func InitLogger(serviceName string) {
	once.Do(
		func() {
			var output zerolog.ConsoleWriter

			output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}
			zerolog.SetGlobalLevel(zerolog.DebugLevel)

			log = zerolog.New(output).
				With().
				Timestamp().
				Str("service", serviceName).
				Logger()
		},
	)
}

// GetLogger возвращает глобальный логгер.
func GetLogger() *zerolog.Logger {
	if log.GetLevel() == zerolog.NoLevel {
		panic("Logger is not initialized. Call InitLogger first.")
	}

	return &log
}

// WithContext добавляет информацию из контекста в логгер.
func WithContext(ctx context.Context) zerolog.Logger {
	logCtx := log.With()

	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		logCtx = logCtx.Str("request_id", requestID)
	}

	return logCtx.Logger()
}

// GetLoggerWithContext возвращает логгер с информацией из контекста.
func GetLoggerWithContext(ctx context.Context) *zerolog.Logger {
	l := WithContext(ctx)
	newLogger := l.With().Logger()

	return &newLogger
}

// LogRequest логирует начало и завершение обработки запроса.
func LogRequest(log zerolog.Logger, ctx context.Context, protocol, method, path string) (context.Context, func()) {
	requestID, _ := ctx.Value("request_id").(string)
	if requestID == "" {
		requestID = "req-" + uuid.New().String()
	}

	ctx = context.WithValue(ctx, "request_id", requestID)

	// Логирование начала обработки запроса
	log.Info().
		Str("request_id", requestID).
		Str("protocol", protocol).
		Str("method", method).
		Str("path", path).
		Msg("Request received")

	start := time.Now()

	// Возвращаем функцию для завершения обработки
	return ctx, func() {
		log.Info().
			Str("request_id", requestID).
			Str("protocol", protocol).
			Str("method", method).
			Str("path", path).
			Int64("duration_ms", time.Since(start).Milliseconds()).
			Msg("Request processed")
	}
}
