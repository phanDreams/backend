package server

import (
	"fmt"
	"net/http"
	"os"
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

	fmt.Printf("DEBUG Start - TLS_ENABLED='%s', TLS_CERT_FILE='%s', TLS_KEY_FILE='%s'\n",
		os.Getenv("TLS_ENABLED"), os.Getenv("TLS_CERT_FILE"), os.Getenv("TLS_KEY_FILE"))
	fmt.Printf("DEBUG TLS startup - certFile='%s' keyFile='%s' enabled=%v\n",
		s.certFile, s.keyFile, s.tlsConfig.Enabled)
	fmt.Printf("DEBUG: TLS_ENABLED=%v, CERT_FILE=%s, KEY_FILE=%s\n",
		s.tlsConfig.Enabled, s.tlsConfig.CertFile, s.tlsConfig.KeyFile)
		
		if s.tlsConfig.Enabled {
			if s.certFile == "" || s.keyFile == "" {
				return fmt.Errorf("TLS is enabled but cert or key file is missing")
			}
			s.logger.Info("Starting HTTPS server...")
			return s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile)
		}
		

	fmt.Println("DEBUG certFile:", s.certFile)
	fmt.Println("DEBUG keyFile:", s.keyFile)
	fmt.Println("DEBUG TLS Enabled:", s.tlsConfig.Enabled)

	s.logger.Info("Starting HTTP server...")
	err := s.httpServer.ListenAndServe()
	return err
}
