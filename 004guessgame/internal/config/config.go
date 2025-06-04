package config

type Config struct {
	DefaultDifficulty	string	`json:"default_difficulty"`
	EnableColors		bool	`json:"enable_colors"`
	EnableSound			bool	`json:"enable_sound"`
	DefaultTimeLimit	int		`json:"default_time_limit"`
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
}

func configPath() (string, error) {
	return "", nil
}