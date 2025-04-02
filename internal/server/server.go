package server

import (
	"net/http"
	"pethelp-backend/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	config     *config.HTTPServerConfig
	logger     *zap.Logger
	httpServer *http.Server
}

func NewHTTPServer(config *config.HTTPServerConfig, logger *zap.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
	}
}

func (s *Server) ListenAndServe(router *gin.Engine) error {
	s.httpServer = &http.Server{
		Addr:         s.config.Address,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
		Handler:      router,
	}
	err := s.httpServer.ListenAndServe()
	return err
}
