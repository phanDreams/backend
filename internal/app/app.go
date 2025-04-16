package app

import (
	"context"
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
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}

	envFilePath := ".env"
	if env := os.Getenv("APP_ENV"); env != "" && env != "local" {
		envFilePath = ""
	}

	if err = config.LoadEnv(envFilePath, logger); err != nil {
		logger.Fatal("Failed to load environment variables", zap.Error(err))
	}

	return fx.New(
		fx.Supply(logger),
		fx.Provide(
			config.NewPostgresConfig,
			config.NewRedisConfig,
			config.LoadHTTPServerConfig,
			postgres.New,
			redis.New,
			server.NewHTTPServer,
			server.NewGinServer,
		),
		// Ensure database is initialized before other components
		fx.Invoke(
			func(s *postgres.Storage, lc fx.Lifecycle) error {
				postgres.ManageLifecycle(s, lc)
				if err := s.Open(context.Background()); err != nil {
					return err
				}
				return nil
			},
			func(r *redis.Storage, lc fx.Lifecycle) error {
				redis.ManageLifecycle(r, lc)
				return nil
			},
		),
		// Add modules after database initialization
		health.Module,
		specialist.Module,
	)
}
