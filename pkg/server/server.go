// Package server запускает сервер и реализует GracefulShutdown.
package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ee-crocush/go-news/pkg/kafka"
	"github.com/ee-crocush/go-news/pkg/logger"
)

var ErrNoServers = errors.New("no servers to start")

// serverTimeout время таймаута сервера.
const serverTimeout = 30 * time.Second

// GracefulServer — интерфейс для серверов, которые поддерживают
// блокирующий запуск и корректное завершение (shutdown).
type GracefulServer interface {
	Start() error
	Shutdown(ctx context.Context) error
}

// ServerManager управляет несколькими серверами.
type ServerManager struct {
	servers []GracefulServer
}

// NewServerManager создает новый менеджер серверов
func NewServerManager(servers ...GracefulServer) *ServerManager {
	return &ServerManager{servers: servers}
}

// StartAll запускает все сервера с graceful shutdown.
func (sm *ServerManager) StartAll(consumer *kafka.Consumer) error {
	if len(sm.servers) == 0 {
		return ErrNoServers
	}

	errChan := make(chan error, len(sm.servers))
	var wg sync.WaitGroup

	// 1️⃣ Стартируем consumer первым
	if consumer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// context без таймаута, чтобы consumer жил до сигнала
			ctxConsumer, cancelConsumer := context.WithCancel(context.Background())
			defer cancelConsumer()

			// Retry для старта consumer на случай, если Kafka ещё не готова
			for {
				err := consumer.Start(ctxConsumer)
				if err != nil && !errors.Is(err, context.Canceled) {
					log := logger.GetLogger()
					log.Err(err).Msg("Kafka consumer failed, retrying in 5s")
					time.Sleep(10 * time.Second)
					continue
				}
				break
			}
		}()
	}

	for _, srv := range sm.servers {
		wg.Add(1)
		go func(s GracefulServer) {
			defer wg.Done()
			if err := s.Start(); err != nil {
				errChan <- err
			}
		}(srv)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		sm.shutdownAll()
		if consumer != nil {
			consumer.Close()
		}
		wg.Wait()
		return err
	case <-sigChan:
		fmt.Println("Received shutdown signal")
		if consumer != nil {
			fmt.Println("Shutting kafka consumer...")
			consumer.Close()
		}
		sm.shutdownAll()
		wg.Wait()
		return nil
	}
}

// shutdownAll корректно завершает все сервера.
func (sm *ServerManager) shutdownAll() {
	var wg sync.WaitGroup
	wg.Add(len(sm.servers))

	for _, srv := range sm.servers {
		go func(s GracefulServer) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := s.Shutdown(ctx); err != nil {
				fmt.Printf("Error shutting down server: %v\n", err)
			}
		}(srv)
	}

	wg.Wait()
}
