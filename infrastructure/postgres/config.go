package postgres

import "users/domain/errors"

type Config struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func NewConfig(
	host string,
	port string,
	database string,
	username string,
	password string,
) (*Config, error) {
	if host == "" {
		return nil, errors.PostgresMissingHost
	}

	if port == "" {
		return nil, errors.PostgresMissingPort
	}

	if database == "" {
		return nil, errors.PostgresMissingDB
	}

	if username == "" {
		return nil, errors.PostgresMissingUser
	}

	if password == "" {
		return nil, errors.PostgresMissingPwd
	}

	return &Config{
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}, nil
}
