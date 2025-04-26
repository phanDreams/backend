package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type TLSConfig struct {
	Enabled  bool
    CertFile string
    KeyFile  string
}

func NewTLSConfig() *TLSConfig {
    enabled := os.Getenv("TLS_ENABLED") == "true"

    return &TLSConfig{
        Enabled:  enabled,
        CertFile: os.Getenv("TLS_CERT_FILE"), // e.g. /home/ubuntu/certs/server.crt
        KeyFile:  os.Getenv("TLS_KEY_FILE"),  // e.g. /home/ubuntu/certs/server.key
    }
}

type PostgresConfig interface {
	DSN() string
}

type RedisConfig interface {
	URI() string
}

type HTTPServerConfig struct {
	Address         string        `yaml:"address"`
	Port            int           `yaml:"port"`
	SecurePort      int           `yaml:"secure_port"`
	Timeout         time.Duration `yaml:"timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type Config struct {
	HTTPServer HTTPServerConfig `yaml:"http_server"`
}

// LoadEnv loads the .env file if the path is provided.
func LoadEnv(path string, logger *zap.Logger) error {
	if path == "" {
		logger.Info("No .env file path provided, using system environment variables")
		return nil
	}

	if err := godotenv.Load(path); err != nil {
		logger.Error("Failed to load .env file", zap.String("path", path), zap.Error(err))
		return err
	}

	logger.Info("Loaded .env file", zap.String("path", path))
	return nil
}

func LoadHTTPServerConfig(logger *zap.Logger) (*HTTPServerConfig, error) {
	yamlConfigPath := "configs/config.yaml"

	var cfg Config

	if err := cleanenv.ReadConfig(yamlConfigPath, &cfg); err != nil {
		logger.Fatal("Error reading config", zap.Error(err))
		return nil, err
	}

	// env := os.Getenv("ENV")
    // if env == "production" {
    //     cfg.HTTPServer.Address = ":443"
    // } else {
    //     cfg.HTTPServer.Address = ":3000"
    // }
	
	// Inject SERVER_ADDRESS from .env if present
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress != "" {
		cfg.HTTPServer.Address = serverAddress
	}
	logger.Info("Loaded HTTP server config", zap.Any("config", cfg.HTTPServer))

	return &cfg.HTTPServer, nil
}