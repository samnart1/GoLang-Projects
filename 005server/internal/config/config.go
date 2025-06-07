package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port 			string			`json:"port"`
	Environment 	string			`json:"environment"`
	LogLevel 		string			`json:"log_level"`
	ReadTimeout 	time.Duration	`json:"read_timeout"`
	WriteTimeout 	time.Duration	`json:"write_timeout"`
	IdleTimeout 	time.Duration	`json:"idle_timeout"`
	ShutdownTimeout time.Duration	`json:"shutdown_timeout"`
	EnableCORS 		bool			`json:"enable_cors"`
	EnableRateLimit bool			`json:"enable_rate_limit"`
	RateLimitRPS 	int				`json:"rate_limit_rps"`
	StaticDir		string			`json:"static_dir"`
	TemplateDir		string			`json:"template_dir"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		ReadTimeout: getDurationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout: getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 30*time.Second),
		EnableCORS: getBoolEnv("ENABLE_CORS", true),
		EnableRateLimit: getBoolEnv("ENABLE_RATE_LIMIT", false),
		RateLimitRPS: getIntEnv("RATE_LIMIT_RPS", 100),
		StaticDir: getEnv("STATIC_DIR", "./web/static"),
		TemplateDir: getEnv("TEMPLATE_DIR", "./web/templates"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	if c.Environment == "" {
		return fmt.Errorf("environment cannot be empty")
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}

	return nil
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}