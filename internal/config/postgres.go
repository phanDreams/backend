package config

import (
	"errors"
	"os"
	"strings"

	"go.uber.org/zap"
)

const (
	dsnEnvName = "PG_DSN"
)

var (
	ErrMissingRequiredEnvVar = errors.New("missing required environment variable")
)

var _ PostgresConfig = (*postgresConfig)(nil)

type postgresConfig struct {
	dsn string
}

// NewPostgresConfig creates a new configuration for PostgresSQL using environment variables.
func NewPostgresConfig(logger *zap.Logger) (PostgresConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if dsn == "" {
		logger.Error("Postgres DSN environment variable is not set", zap.String("env_name", dsnEnvName))
		return nil, ErrMissingRequiredEnvVar
	}

	// Convert direct connection string to session pooler format if needed
	if strings.Contains(dsn, "supabase.co:5432") {
		dsn = strings.Replace(dsn, "supabase.co:5432", "pooler.supabase.com:5432", 1)
		dsn = strings.Replace(dsn, "postgresql://postgres:", "postgres://postgres.", 1)
	}

	logger.Info("Postgres DSN loaded successfully")
	return &postgresConfig{
		dsn: dsn,
	}, nil
}

func (cfg *postgresConfig) DSN() string {
	return cfg.dsn
}