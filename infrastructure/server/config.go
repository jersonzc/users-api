package server

import (
	"time"
	"users/domain/errors"
)

type Config struct {
	Port         int
	Prefix       string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewConfig(
	port int,
	prefix string,
	idleTimeout time.Duration,
	readTimeout time.Duration,
	writeTimeout time.Duration,
) (*Config, error) {
	if port < 0 || port > 65535 {
		return nil, errors.ServerInvalidPort
	}

	if prefix == "" {
		return nil, errors.ServerMissingPrefix
	}

	return &Config{
		Port:         port,
		Prefix:       prefix,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}, nil
}
