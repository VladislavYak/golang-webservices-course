package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Config хранит конфигурацию для подключения к MongoDB.
type Config struct {
	URI        string
	Database   string
	TimeoutSec int
}

// NewMongoClient создает и возвращает новый mongo.Client.
func NewMongoClient(cfg Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeoutSec)*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
