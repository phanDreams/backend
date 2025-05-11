package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RedisDB struct {
	client *redis.Client
}

// Module provides the PostgreSQL database connection pool to the FX container.
var Module = fx.Options(
	fx.Provide(NewClient),
)

// New creates a new Storage with the given config and logger.
func NewClient(lc fx.Lifecycle, logger *zap.Logger) (*RedisDB, error) {
	const operationName = "new_redis_client"
	connString := os.Getenv("REDIS_URI") // Or from a config struct
	if connString == "" {
		getEnvErr := fmt.Errorf("%s: REDIS_URI environment variable not set", operationName)
		logger.Error("failed env get", zap.Error(getEnvErr))
		return nil, getEnvErr
	}

	opt, err := redis.ParseURL(connString)
	if err != nil {
		parseURLErr := fmt.Errorf("%s: failed parse URL: %w", operationName, err)
		logger.Error("failed parse URL", zap.Error(parseURLErr))
		return nil, parseURLErr
	}
	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //For connect
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		pingErr := fmt.Errorf("%s: failed ping Redis database: %w", operationName, err)
		logger.Error("failed ping Redis database", zap.Error(pingErr))
		return nil, pingErr
	}

	logger.Info("Redis connection created successfully")
	redisDB := &RedisDB{client: client}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Redis connection")
			if err := redisDB.client.Close(); err != nil {
				logger.Error("Error closing Redis connection", zap.Error(err))
				return err
			}
			return nil
		},
	})

	return redisDB, nil
}

func (s *RedisDB) Client() *redis.Client {
	return s.client
}
