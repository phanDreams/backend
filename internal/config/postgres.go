package config

import (
	"errors"
	"os"

	"go.uber.org/zap"
)

const (
	dsnEnvName = "PG_DSN"
)

var _ PostgresConfig = (*postgresConfig)(nil)

type postgresConfig struct {
	dsn string
}

// NewPostgresConfig creates a new configuration for PostgresSQL using environment variables.
func NewPostgresConfig(logger *zap.Logger) (PostgresConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if dsn == "" {
		logger.Error("Postgres DSN not found in environment")
		return nil, errors.New("postgres dsn not found")
	}

	logger.Info("Postgres DSN loaded successfully")
	return &postgresConfig{
		dsn: dsn,
	}, nil
}

func (cfg *postgresConfig) DSN() string {
	return cfg.dsn
}