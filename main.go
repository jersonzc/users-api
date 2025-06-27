package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"users/docs"
	"users/infrastructure/dependencies"
	"users/infrastructure/postgres"
	"users/infrastructure/server"
)

// @title           Users API
// @version         1.0
// @description     Interact with user accounts.
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
		if err != nil {
			log.Printf("Error while shutting down otel sdk: %s", err.Error())
		}
	}()

	// Config resources.
	config, err := NewConfig()
	if err != nil {
		log.Printf("Configuration error: %s", err.Error())
		return
	}

	// Swagger
	docs.SwaggerInfo.BasePath = config.Server.Prefix

	// Postgres client.
	postgresClient, err := postgres.NewClient(config.DB)
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
	actions, err := dependencies.NewActions(postgresClient)
	if err != nil {
		log.Printf("Actions error: %s", err.Error())
	}

	// Start HTTP server.
	app := server.Setup(config.Server, actions)
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
}
