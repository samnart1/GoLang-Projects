package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server		ServerConfig	`yaml:"server"`
	Database	DatabaseConfig	`yaml:"database"`	
	Redis		RedisConfig		`yaml:"redis"`
	JWT			JWTConfig		`yaml:"jwt"`
}

type ServerConfig struct {
	Port			string			`yaml:"port"`
	Host			string			`yaml:"host"`
	ReadTimeout		time.Duration	`yaml:"read_timeout"`
	WriteTimeout	time.Duration	`yaml:"write_timeout"`
}

type DatabaseConfig struct {
	Port		string	`yaml:"port"`
	Host		string	`yaml:"host"`
	Username	string	`yaml:"username"`
	Password	string	`yaml:"password"`
	Database	string	`yaml:"database"`
	SSLMode		string	`yaml:"ssl_mode"`
}

type RedisConfig struct {
	Host		string	`yaml:"host"`
	Port		string	`yaml:"port"`
	Password	string	`yaml:"password"`
	DB			int	  	`yaml:"db"`
}

type JWTConfig struct {
	Secret			string			`yaml:"secret"`
	ExpirationTime	time.Duration	`yaml:"expiration_time"`
}

func Load() (*Config, error) {
	config := &Config{
		Server: 	ServerConfig{
			Host: 			getEnv("SERVER_HOST", "localhost"),
			Port: 			getEnv("SERVER_PORT", "8080"),
			ReadTimeout: 	getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: 	getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: 	DatabaseConfig{
			Host: 			getEnv("DB_HOST", "localhost"),
			Port: 			getEnv("DB_PORT", "5432"),
			Password: 		getEnv("DB_PASSWORD", ""),
			Username:		getEnv("DB_USERNAME", "postgres"),
			Database: 		getEnv("DB_NAME", "blog_db"),
			SSLMode: 		getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: 		RedisConfig{
			Host: 			getEnv("REDIS_HOST", "localhost"),
			Port: 			getEnv("REDIS_PORT", "6379"),
			Password: 		getEnv("REDIS_PASSWORD", ""),
			DB: 			getIntEnv("REDIS_DB", 0),
		},
		JWT: 		JWTConfig{
			Secret: 		getEnv("JWT_SECRET", ""),
			ExpirationTime: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		},
	}

	if config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD cannot be empty")
	}

	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET cannot be empty")
	}

	return config, nil
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", 
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.SSLMode,
	)
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func (c *Config) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err != nil {
			return duration
		}
	}
	return defaultValue
}