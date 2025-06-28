package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Log 	LogConfig 		`mapstructure:"log"`
	Server 	ServerConfig	`mapstructure:"server"`
}

type LogConfig struct {
	File 			string	`mapstructure:"file"`
	Level			string	`mapstructure:"level"`
	Format			string	`mapstructure:"format"`
	Output			string	`mapstructure:"output"`
	TimestampFormat	string	`mapstructure:"timestamp_format"`
	MaxSize			int		`mapstructure:"map_size"`
	MaxBackups		int 	`mapstructure:"max_backups"`
	MaxAge			int		`mapstructure:"max_age"`
	Compress		bool	`mapstructure:"compress"`
}

type ServerConfig struct {
	Host string	`mapstructure:"host"`
	Port int	`mapstructure:"port"`
}

func Load() (*Config, error) {
	setDefaults()

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	if cfg.Log.File != "" {
		if err := ensureLogDir(cfg.Log.File); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("log.file", "./logs/simple-logger.log")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.output", "file")
	viper.SetDefault("log.timestamp_format", "2006-01-02 15:04:05")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 7)
	viper.SetDefault("log.compress", true)

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
}

func (c *Config) validate() error {
	validLevels := map[string]bool{
		"debug": true,
		"info": true,
		"warn": true,
		"error": true,
		"fatal": true,
	}

	if !validLevels[c.Log.Level] {
		return fmt.Errorf("invalid log level: %s", c.Log.Level)
	}

	validFormats := map[string]bool{
		"text": true,
		"json": true,
	}

	if !validFormats[c.Log.Format] {
		return fmt.Errorf("invalid log format: %s", c.Log.Format)
	}

	validOutputs := map[string]bool{
		"file":		true,
		"console": 	true,
		"both":		true,
	}

	if !validOutputs[c.Log.Output] {
		return fmt.Errorf("invalid log output: %s", c.Log.Output)
	}

	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	return nil
} 

func ensureLogDir(logFile string) error {
	dir := filepath.Dir(logFile)
	if dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	return nil
}