package app

import (
	"os"
	"time"

	redisStorage "pethelp-backend/internal/database/redis"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"

	apiauth "pethelp-backend/internal/api/auth"
	"pethelp-backend/internal/api/health"
	oauthMod "pethelp-backend/internal/api/oauth"
	"pethelp-backend/internal/config"

	"pethelp-backend/internal/database/postgres"
	"pethelp-backend/internal/logger"
	"pethelp-backend/internal/server"
)

func NewApp() fx.Option {
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}

	envFile := ".env"
	if env := os.Getenv("APP_ENV"); env != "" && env == "local" {
		if err := config.LoadEnv(envFile, logger); err != nil {
			logger.Fatal("failed to load .env", zap.String("envFile", envFile), zap.Error(err))
		}
	}

	return fx.Options(
		// Core providers
		fx.Provide(
			// Logger
			func() *zap.Logger { return logger },
			// Configs
			config.NewPostgresConfig,
			config.NewRedisConfig,
			config.LoadOAuthConf,
			// Redis storage
			redisStorage.New,
			func(s *redisStorage.Storage) *redis.Client {
				return s.Client()
			},

			config.LoadHTTPServerConfig,
			config.NewTLSConfig,
			// Gin engine
			server.NewGinServer,
			// Postgres storage
			postgres.New,
			// HTTP servers
			server.NewHTTPServer,
		),
		// API modules
		health.Module,
		apiauth.Module,
		oauthMod.Module,
		// Server start/stop hooks
		fx.Invoke(
			// Manage postgres storage lifecycle
			postgres.ManageLifecycle,
			redisStorage.ManageLifecycle,
		),
		fx.StartTimeout(20*time.Second),
	)
}
