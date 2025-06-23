package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	postgres *pgxpool.Pool
	dbName   string
	uri      string
	tracer   trace.Tracer
}

func NewClient(config *Config, tracer trace.Tracer) (*Client, error) {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)

	connConfig, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = connPool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &Client{
		postgres: connPool,
		dbName:   config.Database,
		uri:      uri,
		tracer:   tracer,
	}, nil
}
