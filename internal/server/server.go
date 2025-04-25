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
	certFile   string
	keyFile    string
	tlsConfig  *config.TLSConfig
}

func NewHTTPServer(config *config.HTTPServerConfig, logger *zap.Logger, tlsConfig *config.TLSConfig) *Server {
	return &Server{
		config: config,
		logger: logger,
		certFile: tlsConfig.CertFile,
        keyFile:  tlsConfig.KeyFile,
		tlsConfig: tlsConfig,
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

	if s.certFile != "" && s.keyFile != "" && s.tlsConfig.Enabled {
        s.logger.Info("Starting HTTPS server...")
        return s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile)
    }
	s.logger.Info("Starting HTTP server...")
	err := s.httpServer.ListenAndServe()
	return err
}