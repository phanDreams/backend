// import (
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"go.uber.org/fx"
// 	"go.uber.org/zap"

// 	appauth "pethelp-backend/internal/app/auth"
// 	"pethelp-backend/internal/domain/auth"
// 	authinfrastructure "pethelp-backend/internal/infrastructure/auth"
// )

// // ModuleParams holds common dependencies for auth modules
// // It supplies the Gin router, Postgres pool, and logger.
// type ModuleParams struct {
//     fx.In

//     Router *gin.Engine
//     DB     *pgxpool.Pool
//     Logger *zap.Logger
// }

// // AuthModule wires up a /<routePath>/register endpoint for any Registrable entity E
// // using DTO type D. You pass in the Fx module name, DB table name, HTTP route prefix,
// // a constructor for the DTO, and a mapper from DTO to entity.
// func AuthModule[
//     E auth.Registrable,
//     D authinfrastructure.RegistrationDTO,
// ](
//     moduleName string,
//     tableName  string,
//     routePath  string,
//     newDTO     func() D,
//     toEntity   func(D) E,
// ) fx.Option {
//     return fx.Module(moduleName,
//         // Provide the AuthService, along with repository, hasher, and validator
//         fx.Provide(
//             func(p ModuleParams,
//                 hasher auth.PasswordHasher,
//                 validator auth.PasswordValidator,
//             ) *appauth.AuthService {
//                 secret := os.Getenv("JWT_SECRET")
//                 return appauth.NewAuthService(p.DB, p.Logger, hasher, validator, secret, tableName)
//             },
//             // lower-level providers
//             authinfrastructure.NewSQLRepository[E],
//             authinfrastructure.NewFieldsValidator,
//             authinfrastructure.NewBcryptHasher,
//             authinfrastructure.NewPasswordValidator,
//         ),
//         // Invoke the registration handler
//         fx.Invoke(func(p ModuleParams,
//             svc *appauth.AuthService,
//             fv auth.FieldsValidator,
//             repo auth.Repository[E],
//         ) {
//             grp := p.Router.Group(routePath)
//             grp.POST("/register", authinfrastructure.RegisterHandler[E, D](
//                 svc,
//                 fv,
//                 repo,
//                 newDTO,
//                 toEntity,
//                 p.Logger,
//             ))
//             p.Logger.Info("registered auth route", zap.String("route", routePath+"/register"))
//         }),
//     )
// }

package apiauth

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	// "pethelp-backend/internal/infrastructure/postgres"

	appauth "pethelp-backend/internal/app/auth"
	"pethelp-backend/internal/database/postgres"
	"pethelp-backend/internal/domain/auth"
	authinfrastructure "pethelp-backend/internal/infrastructure/auth"
)

// ModuleParams holds common dependencies for auth modules
// It supplies the Gin router, Postgres pool, and logger.
type ModuleParams struct {
    fx.In

    Router *gin.Engine
    DB     *postgres.Storage
    Logger *zap.Logger
}

// AuthModule wires up a /<routePath>/register endpoint for any Registrable entity E
// using DTO type D. You pass in the Fx module name, DB table name, HTTP route prefix,
// a constructor for the DTO, and a mapper from DTO to entity.
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
        // Only use the provided *pgxpool.Pool, Gin router, and Logger
        fx.Provide(
            // Create AuthService with injected DB pool
            func(p ModuleParams,
                hasher auth.PasswordHasher,
                validator auth.PasswordValidator,
                repo auth.Repository[E],
            ) *appauth.AuthService {
                secret := os.Getenv("JWT_SECRET")
                storage := p.DB // Directly use the DB pool if postgres.NewStorage is unavailable
                return appauth.NewAuthService(storage, p.Logger, hasher, validator, secret, tableName)
            },
            // lower-level auth dependencies
            authinfrastructure.NewSQLRepository[E],
            authinfrastructure.NewFieldsValidator,
            authinfrastructure.NewBcryptHasher,
            authinfrastructure.NewPasswordValidator,
        ),
        // Invoke the registration handler
        fx.Invoke(func(p ModuleParams,
            svc *appauth.AuthService,
            fv auth.FieldsValidator,
            repo auth.Repository[E],
        ) {
            grp := p.Router.Group(routePath)
            grp.POST("/register", authinfrastructure.RegisterHandler[E, D](
                svc,
                fv,
                repo,
                newDTO,
                toEntity,
                p.Logger,
            ))
            p.Logger.Info("registered auth route", zap.String("route", routePath+"/register"))
        }),
    )
}