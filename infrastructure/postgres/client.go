package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
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

func (c *Client) Migrate() error {
	db, err := sql.Open("pgx", c.uri)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("failed to close db connection: %v", err)
		}
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

func (c *Client) Modify(ctx context.Context, query string) error {
	tracerCtx, span := c.tracer.Start(ctx, "PostgresClient-Modify")
	defer span.End()

	span.SetAttributes(attribute.String("client.postgres.query", query))

	_, err := c.postgres.Exec(tracerCtx, query)
	if err != nil {
		return fmt.Errorf("unable to modify table: %w", err)
	}

	return nil
}

func (c *Client) Retrieve(ctx context.Context, query string) ([]map[string]string, error) {
	tracerCtx, span := c.tracer.Start(ctx, "PostgresClient-Retrieve")
	defer span.End()

	span.SetAttributes(attribute.String("client.postgres.query", query))

	result, err := c.postgres.Query(tracerCtx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer result.Close()

	var response []map[string]string

	columns := result.FieldDescriptions()
	for result.Next() {
		if result.Err() != nil {
			return nil, fmt.Errorf("postgres: unexpected row scan: %w", result.Err())
		}

		row := make(map[string]string, len(columns))
		values, valuesErr := result.Values()
		if valuesErr != nil {
			return nil, fmt.Errorf("postgres: failed to read row: %w", valuesErr)
		}

		for i, v := range values {
			var value string

			if v != nil {
				switch v.(type) {
				case pgtype.Bool:
					boolVal, _ := v.(pgtype.Bool).BoolValue()
					value = fmt.Sprintf("%v", boolVal)
				case pgtype.Numeric:
					intVal, _ := v.(pgtype.Numeric).Int64Value()
					value = fmt.Sprintf("%d", intVal.Int64)
				case pgtype.Time:
					timeVal, _ := v.(pgtype.Time).TimeValue()
					value = fmt.Sprintf("%v", timeVal)
				default:
					value = fmt.Sprintf("%v", v)
				}

				row[columns[i].Name] = value
			}
		}

		response = append(response, row)
	}

	return response, nil
}
