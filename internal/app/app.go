package app

import (
	"context"
	"os"
	"time"

	redisStorage "pethelp-backend/internal/database/redis"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"

	apiauth "pethelp-backend/internal/api/auth"
	"pethelp-backend/internal/api/health"
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
    if env := os.Getenv("APP_ENV"); env != "" && env != "local" {
        envFile = ""
    }
    if err := config.LoadEnv(envFile, logger); err != nil {
        logger.Fatal("failed to load .env", zap.Error(err))
    }

    return fx.Options(
        // Core providers
        fx.Provide(
            // Logger
            func() *zap.Logger { return logger },
            // Configs
            config.NewPostgresConfig,
            config.NewRedisConfig,
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
        // Server start/stop hooks
        fx.Invoke(
            // Manage postgres storage lifecycle
            postgres.ManageLifecycle,
            redisStorage.ManageLifecycle,
            // Manage HTTP server lifecycle
            func(lc fx.Lifecycle, srv *server.Server, router *gin.Engine) {
                lc.Append(fx.Hook{
                    OnStart: func(ctx context.Context) error {
                        go srv.ListenAndServe(router)
                        return nil
                    },
                    OnStop: func(ctx context.Context) error { return srv.Shutdown(ctx) },
                })
            },
        ),
        fx.StartTimeout(20*time.Second),
    )
}
