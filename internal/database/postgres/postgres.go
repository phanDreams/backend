package postgres

import (
	"context"
	"pethelp-backend/internal/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Storage struct {
	db     *pgxpool.Pool
	cfg    config.PostgresConfig
	logger *zap.Logger
}

// New creates a new Storage with the given config and logger.
func New(pgc config.PostgresConfig, logger *zap.Logger) *Storage {
	return &Storage{
		cfg:    pgc,
		logger: logger,
	}
}

// Open establishes a connection to the PostgreSQL database.
func (s *Storage) Open(ctx context.Context) error {
	config, err := pgxpool.ParseConfig(s.cfg.DSN())
	if err != nil {
		s.logger.Error("failed to parse database config", zap.Error(err))
		return err
	}

	// Add connection pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	
	// Disable statement cache to prevent conflicts
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		s.logger.Error("Failed to open pool connections", zap.Error(err))
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		s.logger.Error("Failed to ping postgres database", zap.Error(err))
		return err
	}

	s.db = pool
	s.logger.Info("Database connection created successfully")
	return nil
}

// Close closes the database connection.
func (s *Storage) Close() {
	if s.db != nil {
		s.db.Close()
		s.logger.Info("Database connection closed")
	}
}

// DB returns the database pool.
func (s *Storage) DB() *pgxpool.Pool {
	return s.db
}

// ManageLifecycle registers Open and Close with the FX lifecycle.
func ManageLifecycle(s *Storage, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.Open(ctx)
		},
		OnStop: func(ctx context.Context) error {
			s.Close()
			return nil
		},
	})
}
