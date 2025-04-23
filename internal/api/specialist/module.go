package specialist

import (
	"context"
	"fmt"
	"os"
	"pethelp-backend/internal/database/postgres"
	"pethelp-backend/internal/domain/service"
	"pethelp-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const specialistRoutePath = "/api/v1/specialists"

type ModuleParams struct {
	fx.In

	Router  *gin.Engine
	Storage *postgres.Storage
	Logger  *zap.Logger
	Lc      fx.Lifecycle
}

var Module = fx.Module("specialist",
	fx.Provide(
		func(p ModuleParams) (*service.AuthService, error) {
			var authService *service.AuthService
			var err error
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return nil, fmt.Errorf("JWT_SECRET environment variable not set")
			}
			p.Lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					secret := os.Getenv("JWT_SECRET")
					if secret == "" {
						return fmt.Errorf("JWT_SECRET environment variable not set")
					}

					db := p.Storage.DB()
					if db == nil {
						return fmt.Errorf("database connection not initialized")
					}

					authService = service.NewAuthService(db, p.Logger, secret)
					return nil
				},
			})
			return authService, err
		},
	),
	fx.Invoke(
		func(p ModuleParams, authService *service.AuthService) {
			specialistGroup := p.Router.Group(specialistRoutePath)
			{
				handler := handlers.RegisterSpecialistHandler(authService, p.Logger)
				specialistGroup.POST("/register", handler)

				p.Logger.Info("Registered specialist routes",
					zap.String("base_path", specialistRoutePath),
					zap.String("register_endpoint", "/register"),
					zap.String("method", "POST"))
			}
		},
	),
)
