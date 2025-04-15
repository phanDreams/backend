package app

import (
	"log"
	"os"
	"pethelp-backend/internal/api/health"
	"pethelp-backend/internal/api/specialist"
	"pethelp-backend/internal/config"
	"pethelp-backend/internal/database/postgres"
	"pethelp-backend/internal/database/redis"
	"pethelp-backend/internal/logger"
	"pethelp-backend/internal/server"
	"pethelp-backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewApp() *fx.App {
	envFilePath := ".env"
	if env := os.Getenv("APP_ENV"); env != "" && env != "local" {
		envFilePath = ""
	}

	if err := config.LoadEnv(envFilePath, nil); err != nil {
		log.Fatal("Failed to load environment variables:", err)
	}

	return fx.New(
		fx.Provide(
			// Core dependencies
			logger.New,
			config.NewPostgresConfig,
			config.NewRedisConfig,
			config.LoadHTTPServerConfig,
			
			// Server components
			server.NewGinRouter,
			server.NewHTTPServer,
			
			// Database
			postgres.New,
			redis.New,
			
			// Services
			func(db *postgres.Storage, logger *zap.Logger) *service.AuthService {
				return service.NewAuthService(db.DB(), logger, "your-jwt-secret-here")
			},
		),
		
		// Modules
		specialist.Module,
		health.Module,

		fx.Invoke(
			// This will ensure the server starts
			func(
				router *gin.Engine,
				server *server.Server,
				logger *zap.Logger,
				config *config.HTTPServerConfig,
			) {
				go func() {
					logger.Info("Starting HTTP server", zap.String("address", config.Address))
					if err := server.ListenAndServe(router); err != nil {
						logger.Fatal("Server failed to start", zap.Error(err))
					}
				}()
			},
		),
	)
}
