package specialist

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"pethelp-backend/internal/handlers"
	"pethelp-backend/internal/service"
)

const (
	specialistRoutePath = "/api/v1/specialists"
)

var Module = fx.Module("specialist",
	fx.Provide(
		// Add any specialist-specific providers here
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(
	router *gin.Engine,
	authService *service.AuthService,
	logger *zap.Logger,
) {
	if router == nil {
		logger.Fatal("Router is nil in registerRoutes")
		return
	}

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
