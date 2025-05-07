package server

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewGinServer(lc fx.Lifecycle, logger *zap.Logger, server *Server) *gin.Engine {
	router := gin.Default()

	// Register custom validators
	if _, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Add any custom validators here if needed
		logger.Info("Validator engine initialized")
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Listening on", zap.String("address", server.config.Address))
	
			env := os.Getenv("ENV")
			if env == "production" {
				return server.ListenAndServe(router)
			} else {
				go func() {
					if err := server.ListenAndServe(router); err != nil && err != http.ErrServerClosed {
						logger.Fatal("Failed to start server", zap.Error(err))
					}
				}()
				return nil
			}
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server...")
			return server.httpServer.Shutdown(ctx)
		},
	})
	
	return router
}