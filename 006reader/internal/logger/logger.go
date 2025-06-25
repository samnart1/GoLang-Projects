package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func New() (*Logger, error) {
	config := zap.NewProductionConfig()

	if os.Getenv("ENVIRONMENT") == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	level := os.Getenv("LOG_LEVEL")
	if level != "" {
		parsedLevel, err := zapcore.ParseLevel(level)
		if err == nil {
			config.Level.SetLevel(parsedLevel)
		}
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}

func String(key, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func Duration(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}