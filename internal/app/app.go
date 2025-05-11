package app

import (
	"os"
	"pethelp-backend/internal/api/health"
	"pethelp-backend/internal/api/oauth"
	"pethelp-backend/internal/api/specialist"
	"pethelp-backend/internal/config"
	"pethelp-backend/internal/database/postgres"
	"pethelp-backend/internal/database/redis"
	"pethelp-backend/internal/logger"
	"pethelp-backend/internal/server"
	"time"

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
		// Core services
		fx.Provide(
			config.NewPostgresConfig,
			config.NewRedisConfig,
			config.NewServersConfig,
			config.NewTLSConfig,
			config.LoadOAuthConf,
			// postgres.New,
			// redis.New,
			server.NewHTTPServer,
			server.NewGinServer,
		),

		// API modules
		health.Module,
		postgres.Module,
		redis.Module,
		specialist.Module,
		oauth.Module,

		fx.StartTimeout(10*time.Second),
	)
}
