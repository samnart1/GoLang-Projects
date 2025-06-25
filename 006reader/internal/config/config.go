package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Reader		ReaderConfig	`mapstructure:"reader"`
	Formatter	FormatterConfig	`mapstructure:"formatter"`
	Logger		LoggerConfig	`mapstructure:"logger"`
}

type ReaderConfig struct {
	BufferSize	int		`mapstructure:"buffer_size"`
	MaxFileSize	int64	`mapstructure:"max_file_size"`
	Encoding	string	`mapstructure:"encoding"`
}

type FormatterConfig struct {
	MaxWidth 	int		`mapstructure:"max_width"`
	Theme		string	`mapstructure:"theme"`
}

type LoggerConfig struct {
	Level	string	`mapstructure:"level"`
	Format	string	`mapstructure:"format"`
	Output	string	`mapstructure:"output"`
}

func Load() (*Config, error) {

	SetDefaults()

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validate(config); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return config, nil
}

func SetDefaults() {
	viper.SetDefault("reader.buffer_size", 8192)
	viper.SetDefault("reader.max_file_size", 100*1024*1024)
	viper.SetDefault("reader.encoding", "utf-8")

	viper.SetDefault("formatter.max_width", 120)
	viper.SetDefault("formatter.theme", "default")

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")
	viper.SetDefault("logger.output", "stdout")
}

func validate(config *Config) error {
	if config.Reader.BufferSize <= 0 {
		return fmt.Errorf("reader.buffer_size must be positive")
	}

	if config.Reader.MaxFileSize <= 0 {
		return fmt.Errorf("reader.max_file_size must be positive")
	}

	if config.Formatter.MaxWidth <= 0 {
		return fmt.Errorf("formatter.max_width must be positive")
	}

	return nil
}