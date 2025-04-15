package server

import (
	"context"
	"net/http"
	"pethelp-backend/internal/config"
	"pethelp-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewGinRouter(logger *zap.Logger) *gin.Engine {
    gin.SetMode(gin.ReleaseMode)
    router := gin.New()

    // Add basic middleware
    router.Use(gin.Recovery())
    router.Use(middleware.ZapLogger(logger))

    // Debug: Log that the router was created
    logger.Info("Gin router initialized")

    return router
}

func NewHTTPServerHandler(
    router *gin.Engine,
    logger *zap.Logger,
    lc fx.Lifecycle,
    config *config.HTTPServerConfig,
) *http.Server {
    srv := &http.Server{
        Addr:    config.Address,
        Handler: router,
    }

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            logger.Info("Starting HTTP server", 
                zap.String("address", config.Address))
            go func() {
                if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                    logger.Fatal("Failed to start server", zap.Error(err))
                }
            }()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return srv.Shutdown(ctx)
        },
    })

    return srv
}
