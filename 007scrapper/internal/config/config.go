package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Timeout		time.Duration
	UserAgent	string

	ServerPort	string
	ServerHost	string

	MaxConcurrent	int
	RetryAttempts	int
	RateLimit		time.Duration
}

func Load() *Config {
	cfg := &Config{
		Timeout: 30 * time.Second,
		UserAgent: "Go-Web-Scraper/1.0",
		ServerPort: "8080",
		ServerHost: "localhost",
		MaxConcurrent: 10,
		RetryAttempts: 3,
		RateLimit: 100 * time.Millisecond,
	}

	if timeout := os.Getenv("SCRAPER_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			cfg.Timeout = d
		}
	}

	if userAgent := os.Getenv("SCRAPER_USER_AGENT"); userAgent != "" {
		cfg.UserAgent = userAgent
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.ServerPort = port
	}

	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.ServerHost = host
	}

	if maxConcurrent := os.Getenv("MAX_CONCURRENT"); maxConcurrent != "" {
		if n, err := strconv.Atoi(maxConcurrent); err == nil && n > 0 {
			cfg.MaxConcurrent = n
		}
	}

	if retryAttempts := os.Getenv("RETRY_ATTEMPS"); retryAttempts != "" {
		if n, err := strconv.Atoi(retryAttempts); err == nil {
			cfg.RetryAttempts = n
		}
	}

	if rateLimit := os.Getenv("RATE_LIMIT"); rateLimit != "" {
		if d, err := time.ParseDuration(rateLimit); err == nil {
			cfg.RateLimit = d
		}
	}

	return cfg
}

func (c *Config) ServerAddress() string {
	return c.ServerHost + ":" + c.ServerPort
}