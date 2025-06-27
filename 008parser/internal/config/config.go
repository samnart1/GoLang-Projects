package config

import "github.com/spf13/viper"

type Config struct {
	DefaultFormat	string	`mapstructure:"DEFAULT_FORMAT"`
	DefaultIndent	int		`mapstructure:"DEFAULT_INDENT"`
	EnableColors	bool	`mapstructure:"ENABLE_COLORS"`

	ServerHost		string	`mapstructure:"SERVER_HOST"`
	ServerPort		int		`mapstructure:"SERVER_PORT"`

	MaxFileSize		int64	`mapstructure:"MAX_FILE_SIZE"`
	HTTPTimeout		int		`mapstructure:"HTTP_TIMEOUT"`

	NoColor	bool
	Verbose	bool
	Quiet	bool
	Debug	bool
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	cfg.NoColor = viper.GetBool("no-color")
	cfg.Verbose = viper.GetBool("verbose")
	cfg.Quiet = viper.GetBool("quiet")
	cfg.Debug = viper.GetBool("debug")

	return cfg, nil
}