package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"pethelp-backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// DB is the PostgreSQL database connection pool.
type DB struct {
	pool *pgxpool.Pool
}

// Module provides the PostgreSQL database connection pool to the FX container.
var Module = fx.Options(
	fx.Provide(NewPGPool),
)

// NewStorage creates a new database storage.
func NewPGPool(lc fx.Lifecycle, dbConf *config.Config, logger *zap.Logger) (*DB, error) {
	const operationName = "NewPGPool"
	connString := os.Getenv("PG_DSN") // Or from a config struct
	if connString == "" {
		getEnvErr := fmt.Errorf("%s: PG_DSN environment variable not set", operationName)
		logger.Error("failed env get", zap.Error(getEnvErr))
		return nil, getEnvErr
	}

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		getEnvErr := fmt.Errorf("%s: failed to parse connection string: %w", operationName, err)
		logger.Error("failed env get", zap.Error(getEnvErr))
		return nil, getEnvErr
	}

	poolConfig.MaxConns = dbConf.PostgresDB.MaxPoolSize
	poolConfig.MinConns = dbConf.PostgresDB.MinPoolSize
	poolConfig.MaxConnLifetime = dbConf.PostgresDB.MaxLifetime
	poolConfig.MaxConnIdleTime = dbConf.PostgresDB.IdleTimeout
	poolConfig.ConnConfig.ConnectTimeout = dbConf.PostgresDB.ConnectionTimeout

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //For connect
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to database: %w", operationName, err)
	}

	// Add a lifecycle hook to close the connection pool.
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection pool")
			dbpool.Close()
			return nil
		},
	})

	// Test the connection pool
	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to connect to database: %w", operationName, err)
	}

	log.Println("Connected to PostgreSQL pool")
	return &DB{dbpool}, nil

}

// DB returns the database connection pool.
func (s *DB) Pool() *pgxpool.Pool {
	return s.pool
}
