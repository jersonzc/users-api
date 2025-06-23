package main

import (
	"os"
	"strconv"
	"time"
	"users/infrastructure/postgres"
	"users/infrastructure/server"
)

type Config struct {
	Server *server.Config
	DB     *postgres.Config
}

func NewConfig() (*Config, error) {
	serverConfig, err := server.NewConfig(
		getInt("API_PORT", 8080),
		get("PREFIX", "/app"),
		getDuration("SERVER_READ_TIMEOUT", 10, time.Second),
		getDuration("SERVER_WRITE_TIMEOUT", 10, time.Second),
	)
	if err != nil {
		return nil, err
	}

	dbConfig, err := postgres.NewConfig(
		get("DB_HOST", "localhost"),
		get("DB_PORT", "5432"),
		get("DB_NAME", "users"),
		get("DB_USER", "postgres"),
		get("DB_PASSWORD", "postgres"),
	)
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: serverConfig,
		DB:     dbConfig,
	}, nil
}

func get(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if valueToInt, err := strconv.Atoi(value); err == nil {
			return valueToInt
		}
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration, unit time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if valueToInt, err := strconv.Atoi(value); err == nil {
			return time.Duration(valueToInt) * unit
		}
	}
	return defaultValue * unit
}
