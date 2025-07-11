package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	pool    *pgxpool.Pool
	queries *Queries
	dbName  string
	uri     string
	tracer  trace.Tracer
}

func NewClient(config *Config) (*Client, error) {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)

	connConfig, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	connPool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = connPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &Client{
		pool:    connPool,
		queries: New(connPool),
		dbName:  config.Database,
		uri:     uri,
		tracer:  otel.Tracer("PostgresClient"),
	}, nil
}

func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) Migrate() error {
	db, err := sql.Open("pgx", c.uri)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		err = errors.Join(err, db.Close())
	}()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance("file://migrations/", c.dbName, driver)
	if err != nil {
		return err
	}

	if err = migration.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
