package app

import (
	"os"
	"pethelp-backend/internal/api/health"
	"pethelp-backend/internal/api/specialist"
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

	// Load environment variables before creating the FX app
	logger, err := logger.New()
	if err = config.LoadEnv(envFilePath, logger); err != nil {
		logger.Fatal("Failed to load environment variables", zap.Error(err))
	}

	return fx.New(
		fx.Supply(logger), // Supply the already created logger
		health.Module,
		specialist.Module,
		fx.Provide(
			config.NewPostgresConfig,
			config.NewRedisConfig,
			config.LoadHTTPServerConfig,
			postgres.New,
			redis.New,
			server.NewHTTPServer,
			server.NewGinServer,
		),
		fx.Invoke(
			func(s *postgres.Storage, lc fx.Lifecycle) {
				postgres.ManageLifecycle(s, lc)
			},
			func(r *redis.Storage, lc fx.Lifecycle) {
				redis.ManageLifecycle(r, lc)
			},
		),
	)
}
