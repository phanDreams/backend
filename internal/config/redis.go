package config

import (
	"errors"
	"os"

	"go.uber.org/zap"
)

const (
	redisURLEnvName = "REDIS_URI"
)

var _ RedisConfig = (*redisConfig)(nil)

type redisConfig struct {
	uri string
}

// NewRedisConfig creates a new configuration for Redis using environment variables
func NewRedisConfig(logger *zap.Logger) (RedisConfig, error) {
	uri := os.Getenv(redisURLEnvName)
	if uri == "" {
		logger.Error("Redis URI not found in environment")
		return nil, errors.New("redis uri not found")
	}

	logger.Info("Redis URI loaded successfully")
	return &redisConfig{
		uri: uri,
	}, nil
}

// URI returns the Redis URI
func (cfg *redisConfig) URI() string {
	return cfg.uri
}
