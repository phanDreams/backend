package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates and returns zap logger with configured log level
func New() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stdout"}

	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		LineEnding:  zapcore.DefaultLineEnding,
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		if os.Getenv("APP_ENV") == "local" || (len(os.Args) > 1 && os.Args[1] == ".") {
			logLevel = "debug"
		} else {
			logLevel = "info"
		}
	}

	level, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return nil, err
	}
	cfg.Level = level

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
