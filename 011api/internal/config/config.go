package config

type Config struct {
	Environment string
	// Server		ServerConfig
}

type ServerConfig struct {}

func LoadConfig(Environment string) (*Config, error) {
	return &Config{
		Environment: "Samuel",
	}, nil
}