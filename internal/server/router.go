package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewGinServer(lc fx.Lifecycle, logger *zap.Logger, server *Server) *gin.Engine {
	router := gin.Default()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Listening on", zap.String("address", server.config.Address))
				if err := server.ListenAndServe(router); err != nil {
					logger.Fatal("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server...")
			return server.httpServer.Shutdown(ctx)
		},
	})

	return router
}
