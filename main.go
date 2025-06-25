package main

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"log"
	"os"
	"os/signal"
	"users/infrastructure/dependencies"
	"users/infrastructure/postgres"
	"users/infrastructure/server"
)

func main() {
	log.Println("Starting service")

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Create tracer
	tracer := otel.Tracer("app-server")

	// Config resources.
	config, err := NewConfig()
	if err != nil {
		log.Printf("Configuration error: %s", err.Error())
		return
	}

	// Postgres client.
	postgresClient, err := postgres.NewClient(config.DB, tracer)
	if err != nil {
		log.Printf("Database error: %s", err.Error())
		return
	}

	// Run migrations
	err = postgresClient.Migrate()
	if err != nil {
		log.Printf("Database migration error: %s", err.Error())
		return
	}

	// Link actions
	actions, err := dependencies.NewActions(postgresClient, tracer)
	if err != nil {
		log.Fatalf("Actions error: %s", err.Error())
	}

	// Start HTTP server.
	app := server.Setup(config.Server, actions, tracer)
	appErr := make(chan error, 1)
	go func() {
		appErr <- app.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-appErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = app.Shutdown(context.Background())
	return
}
