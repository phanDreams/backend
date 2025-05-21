package oauth

import (
	tokenSrv "pethelp-backend/internal/app/oauth/service"
	tokenEnt "pethelp-backend/internal/domain/oauth"
	"pethelp-backend/internal/infrastructure/oauth/handlers"
	tokenRepo "pethelp-backend/internal/infrastructure/oauth/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const oauthCoreRoutePath = "/api/v1/oauth"

type ModuleParams struct {
	fx.In

	Router *gin.Engine
	Logger *zap.Logger
}

var Module = fx.Module("google_oauth",
	fx.Provide(
		fx.Annotate(
			tokenRepo.NewOAuthTokenRepo,
			fx.As(new(tokenEnt.OAuthTokenRepository)),
		),

		fx.Annotate(
			tokenSrv.NewOAuthService,
			fx.As(new(tokenEnt.OAuthTokenService)),
		),

		handlers.NewOAuthHandlers,
	),
	fx.Invoke(
		func(mp ModuleParams, handler *handlers.OAuthHandlers) {
			oauthGroup := mp.Router.Group(oauthCoreRoutePath)
			{
				oauthGroup.GET("/google", handler.SignInWithProvider)
				oauthGroup.GET("/google/callback", handler.ProviderCallback)

				mp.Logger.Info("Registered google OAuth routes",
					zap.String("base_path", oauthCoreRoutePath),
					zap.String("register_endpoint", "/google"),
					zap.String("method", "GET"))
			}
		},
	),
)
