package specialist

import (
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
		func(s *postgres.Storage, logger *zap.Logger) *service.AuthService {
			return service.NewAuthService(s.DB(), logger, os.Getenv("JWT_SECRET"))
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
