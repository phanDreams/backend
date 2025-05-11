package oauth

import (
	"fmt"

	"pethelp-backend/internal/config"
	"pethelp-backend/internal/database/redis"
	"pethelp-backend/internal/domain/service"
	"pethelp-backend/internal/handlers"
	"pethelp-backend/internal/repository"

	"github.com/gin-gonic/gin" //  Ensure you are using pgxpool
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const oauthCoreRoutePath = "/api/v1/oauth"
const operationName = "oauth_api_module"

type ModuleParams struct {
	fx.In

	Router *gin.Engine
	Redis  *redis.RedisDB
	Logger *zap.Logger
	Lc     fx.Lifecycle
}

var Module = fx.Module("googlea_oauth",
	fx.Provide(
		func(mp ModuleParams) (*service.OAuthUserService, error) {
			redisClient := mp.Redis.Client()
			if redisClient == nil {
				return nil, fmt.Errorf("%s redis database connection not initialized", operationName)
			}

			googleConf, err := config.LoadOAuthConf()
			if err != nil {
				return nil, fmt.Errorf("%s failed to load OAuth configuration: %w", operationName, err)
			}

			oauthRepo := repository.NewOAuthUserRepository(redisClient)
			oauthService := service.NewOAuthService(oauthRepo, googleConf)

			return oauthService, nil
		},
	),
	fx.Invoke(
		func(mp ModuleParams, oauthService *service.OAuthUserService) {
			specialistGroup := mp.Router.Group(oauthCoreRoutePath)
			{
				handler := handlers.NewOAuthHandlers(oauthService)

				specialistGroup.GET("/google", handler.SignInWithProvider)
				specialistGroup.GET("/google/callback", handler.ProviderCallbackHandler)
				specialistGroup.GET("/success", handler.Success)

				mp.Logger.Info("Registered google OAuth routes",
					zap.String("base_path", oauthCoreRoutePath),
					zap.String("register_endpoint", "/google"),
					zap.String("method", "GET"))
			}
		},
	),
)
