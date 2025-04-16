package specialist

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"pethelp-backend/internal/database/postgres"
	"pethelp-backend/internal/domain/service"
	"pethelp-backend/internal/handlers"
)

const (
	specialistRoutePath = "/api/v1/specialists"
)

var Module = fx.Module("specialist",
	fx.Provide(
		func(storage *postgres.Storage, lc fx.Lifecycle, logger *zap.Logger) (*service.AuthService, error) {
			// Ensure database is initialized before creating AuthService
			if err := storage.Open(context.Background()); err != nil {
				logger.Fatal("Failed to initialize database connection", zap.Error(err))
				return nil, err
			}
			
			return service.NewAuthService(storage.DB(), logger, os.Getenv("JWT_SECRET")), nil
		},
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(
	router *gin.Engine,
	authService *service.AuthService,
	logger *zap.Logger,
) {

	specialistGroup := router.Group(specialistRoutePath)
	{
		handler := handlers.RegisterSpecialistHandler(authService, logger)
		specialistGroup.POST("/register", handler)

		logger.Info("Registered specialist routes",
			zap.String("base_path", specialistRoutePath),
			zap.String("register_endpoint", "/register"),
			zap.String("method", "POST"))
	}
}
