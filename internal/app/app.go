package app

import (
	"os"
	"pethelp-backend/internal/api/health"
	"pethelp-backend/internal/config"
	"pethelp-backend/internal/database/postgres"
	"pethelp-backend/internal/database/redis"
	"pethelp-backend/internal/logger"
	"pethelp-backend/internal/server"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewApp() *fx.App {
	envFilePath := ".env"
	if env := os.Getenv("APP_ENV"); env != "" && env != "local" {
		envFilePath = ""
	}

	return fx.New(
		health.Module,
		fx.Provide(
			logger.New,
			config.NewPostgresConfig,
			config.NewRedisConfig,
			config.LoadHTTPServerConfig,
			postgres.New,
			redis.New,
			server.NewHTTPServer,
			server.NewGinServer,
		),
		fx.Invoke(
			func(logger *zap.Logger) error {
				return config.LoadEnv(envFilePath, logger)
			},
			func(s *postgres.Storage, lc fx.Lifecycle) {
				postgres.ManageLifecycle(s, lc)
			},
			func(r *redis.Storage, lc fx.Lifecycle) {
				redis.ManageLifecycle(r, lc)
			},
		),
	)
}