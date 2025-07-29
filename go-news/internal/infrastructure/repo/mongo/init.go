// Package mongo содержит реализацию репозиториев для работы с MongoDB.
package mongo

import (
	"GoNews/internal/infrastructure/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Init инициализирует БД Монго и возвращает клиент MongoDB и выбранную базу.
func Init(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	uri := cfg.MongoDB.URI().String()
	clientOpts := options.Client().ApplyURI(uri).SetConnectTimeout(cfg.MongoDB.ConnectTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("Init.MongoDB.Connect: %w", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, fmt.Errorf("Init.MongoDB.Ping: %w", err)
	}

	db := client.Database(cfg.MongoDB.Database)
	return client, db, nil
}
