package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type AuthConfig struct {
    Table          string        `envconfig:"AUTH_TABLE"          default:"specialists"`
    JWTSecret      []byte        `envconfig:"JWT_SECRET"          required:"true"`
    DefaultTimeout time.Duration `envconfig:"AUTH_TIMEOUT"        default:"5s"`
    AccessTTL      time.Duration `envconfig:"ACCESS_TOKEN_TTL"    default:"15m"`
    RefreshTTL     time.Duration `envconfig:"REFRESH_TOKEN_TTL"   default:"168h"`
}


func LoadAuthConfig(logger *zap.Logger) (AuthConfig, error) {
	var cfg AuthConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("loading auth config: %w", err)
	}
	return cfg, nil
}