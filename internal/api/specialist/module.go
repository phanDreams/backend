package specialist

import (
	"fmt"
	"os"

	"pethelp-backend/internal/database/postgres"
	"pethelp-backend/internal/domain/service"
	"pethelp-backend/internal/handlers"

	"github.com/gin-gonic/gin" //  Ensure you are using pgxpool
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const specialistRoutePath = "/api/v1/specialists"

type ModuleParams struct {
	fx.In

	Router *gin.Engine
	DB     *postgres.DB
	Logger *zap.Logger
	Lc     fx.Lifecycle
}

var Module = fx.Module("specialist",
	fx.Provide(
		func(p ModuleParams) (*service.AuthService, error) {

			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return nil, fmt.Errorf("JWT_SECRET environment variable not set")
			}

			dbPool := p.DB.Pool() // Get the db connection
			if dbPool == nil {
				return nil, fmt.Errorf("database connection not initialized")
			}
			authService, err := service.NewAuthService(dbPool, p.Logger, secret) // Create AuthService
			if err != nil {
				return nil, fmt.Errorf("failed to create AuthService: %w", err)
			}

			return authService, err

		},
	),
	fx.Invoke(
		func(p ModuleParams, authService *service.AuthService) {
			specialistGroup := p.Router.Group(specialistRoutePath)
			{
				handler := handlers.RegisterSpecialistHandler(authService, p.Logger) // Get Handler

				specialistGroup.POST("/register", handler)

				p.Logger.Info("Registered specialist routes",
					zap.String("base_path", specialistRoutePath),
					zap.String("register_endpoint", "/register"),
					zap.String("method", "POST"))
			}
		},
	),
)
