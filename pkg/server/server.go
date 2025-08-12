// Package server запускает сервер и реализует GracefulShutdown.
package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ee-crocush/go-news/pkg/kafka"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), serverTimeout)
	defer cancel()

	errChan := make(chan error, len(sm.servers))
	var wg sync.WaitGroup

	for _, srv := range sm.servers {
		wg.Add(1)
		go func(s GracefulServer) {
			defer wg.Done()
			if err := s.Start(); err != nil {
				errChan <- err
			}
		}(srv)
	}

	// Не забываем про кафку
	if consumer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxConsumer, cancelConsumer := context.WithCancel(context.Background())
			defer cancelConsumer()
			if err := consumer.Start(ctxConsumer); err != nil {
				errChan <- err
			}
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		sm.shutdownAll(ctx)
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
		sm.shutdownAll(ctx)
		wg.Wait()
		return nil
	}
}

// shutdownAll корректно завершает все сервера.
func (sm *ServerManager) shutdownAll(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(sm.servers))

	for _, srv := range sm.servers {
		go func(s GracefulServer) {
			defer wg.Done()
			if err := s.Shutdown(ctx); err != nil {
				fmt.Printf("Error shutting down server: %v\n", err)
			}
		}(srv)
	}

	wg.Wait()
}
