package apiauth

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"

	appauth "pethelp-backend/internal/app/auth"
	"pethelp-backend/internal/config"
	"pethelp-backend/internal/database/postgres"
	"pethelp-backend/internal/domain/auth"
	authinfrastructure "pethelp-backend/internal/infrastructure/auth"
)

// ModuleParams holds common dependencies for auth modules.
// It supplies the Gin router, Postgres pool, Logger, and Redis client.
type ModuleParams struct {
	fx.In

	Router *gin.Engine
	DB     *postgres.Storage
	Cache  *redis.Client
	Logger *zap.Logger
}

// AuthModule wires up /<routePath> endpoints for registration, login, and refresh.
// E is the entity type, D is its DTO type.
func AuthModule[
	E auth.Registrable,
	D authinfrastructure.RegistrationDTO,
](
	moduleName string,
	tableName  string,
	routePath  string,
	newDTO     func() D,
	toEntity   func(D) E,
) fx.Option {
	return fx.Module(moduleName,
		fx.Provide(
			config.LoadAuthConfig,
			authinfrastructure.NewSQLRepository[E],
			authinfrastructure.NewFieldsValidator,
			authinfrastructure.NewBcryptHasher,
			authinfrastructure.NewPasswordValidator,
		),
		fx.Supply(tableName),
		fx.Provide(func(
				p ModuleParams,
				cfg config.AuthConfig,
				hasher auth.PasswordHasher,
				validator auth.PasswordValidator,
				repo auth.Repository[E],
				cache appauth.RedisClient,
				// tbl string,
			) *appauth.AuthService {
				return appauth.NewAuthService(
					p.DB,
					cfg,
					p.Logger,
					hasher,
					validator,
					cache,
					// tbl,
				)
			},
		),
	   

		// Mount handlers
		fx.Invoke(func(p ModuleParams, svc *appauth.AuthService, fv auth.FieldsValidator, repo auth.Repository[E]) {
			grp := p.Router.Group(routePath)
			grp.POST("/register", authinfrastructure.RegisterHandler[E, D](
				svc,
                fv,
                repo,
                newDTO,
                toEntity,
                p.Logger,
			))
			grp.POST("/login", authinfrastructure.LoginHandler(
				svc,
				p.Logger,
			))
			grp.POST("/refresh", authinfrastructure.RefreshHandler(
				svc,
				p.Logger,
			))
			p.Logger.Info("registered auth routes", zap.String("base", routePath))
		}),
	)
}
