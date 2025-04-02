package redis

import (
	"context"
	"pethelp-backend/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Storage struct {
	client *redis.Client
	cfg    config.RedisConfig
	logger *zap.Logger
}

// New creates a new Storage with the given config and logger.
func New(rc config.RedisConfig, logger *zap.Logger) *Storage {
	return &Storage{
		cfg:    rc,
		logger: logger,
	}
}

// Open establishes a connection to the Redis database.
func (s *Storage) Open(ctx context.Context) error {
	opt, _ := redis.ParseURL(s.cfg.URI())
	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		s.logger.Error("Failed to ping Redis database", zap.Error(err))
		return err
	}

	s.client = client
	s.logger.Info("Redis connection created successfully")
	return nil
}

// Close closes the Redis connection.
func (s *Storage) Close() {
	if s.client != nil {
		s.client.Close()
		s.logger.Info("Redis connection closed")
	}
}

// Client returns the Redis client.
func (s *Storage) Client() *redis.Client {
	return s.client
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
