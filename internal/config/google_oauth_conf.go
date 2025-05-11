package config

import (
	"fmt"
	"os"

	"pethelp-backend/internal/utils"
)

type GoogleOAuthConfig struct {
	ClientID          string
	ClientSecret      string
	ClientCallbackURL string
	SessionSecret     string
}

func LoadOAuthConf() (GoogleOAuthConfig, error) {

	operationName := "load_oauth_config"
	cfg := GoogleOAuthConfig{}

	secret, err := utils.GenerateRandomKey(32)
	if err != nil {
		return cfg, fmt.Errorf("%s: failed to generate session secret: %w", operationName, err)
	}

	err = os.Setenv("SESSION_SECRET", secret)
	if err != nil {
		return cfg, fmt.Errorf("%s: failed to set session secret env: %w", operationName, err)
	}

	// err = godotenv.Load(".env")
	// if err != nil && os.Getenv("DEPLOYMENT_ENVIRONMENT") != "production" {
	// 	return cfg, fmt.Errorf("%s: failed to load env file to environment", operationName)
	// }

	cfg = GoogleOAuthConfig{
		ClientID:          os.Getenv("CLIENT_ID"),
		ClientSecret:      os.Getenv("CLIENT_SECRET"),
		ClientCallbackURL: os.Getenv("CLIENT_CALLBACK_URL"),
		SessionSecret:     os.Getenv("SESSION_SECRET"),
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.ClientCallbackURL == "" || cfg.SessionSecret == "" {
		return cfg, fmt.Errorf("%s: some env variables are missing", operationName)
	}

	return cfg, nil
}
