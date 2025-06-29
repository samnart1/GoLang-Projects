package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Default 		DefaultConfig 		`mapstructure:"default"`
	Sound			SoundConfig			`mapstructure:"sound"`
	Display			DisplayConfig		`mapstructure:"display"`
	Notifications	NotificationConfig	`mapstructure:"notifications"`
}

type DefaultConfig struct {
	Sound					bool	`mapstructure:"sound"`
	DesktopNotifications	bool	`mapstructure:"desktop_notifications"`
	Color					bool	`mapstructure:"color"`
}

type SoundConfig struct {
	Enabled	bool	`mapstructure:"enabled"`
	File	string	`mapstructure:"file"`
	Volume	int	`mapstructure:"volume"`
}

type DisplayConfig struct {
	ShowMilliseconds	bool	`mapstructure:"show_milliseconds"`
	Animation			bool	`mapstructure:"animation"`
	ColorSheme			string	`mapstructure:"color_scheme"`
	Color				bool	`mapstructure:"color"`
}

type NotificationConfig struct {
	Desktop  bool `mapstructure:"desktop"`
	Terminal bool `mapstructure:"terminal"`
}

func Load() (*Config, error) {
	config := &Config{
		Default: DefaultConfig{
			Sound: true,
			DesktopNotifications: true,
			Color: true,
		},
		Sound: SoundConfig{
			Enabled: true,
			File: getDefaultSoundFile(),
			Volume: 80,
		},
		Display: DisplayConfig{
			ShowMilliseconds: false,
			Animation: true,
			ColorSheme: "default",
			Color: true,
		},
		Notifications: NotificationConfig{
			Desktop: true,
			Terminal: true,
		},
	}

	viper.SetDefault("default.sound", config.Default.Sound)
	viper.SetDefault("default.desktop_notifications", config.Default.DesktopNotifications)
	viper.SetDefault("default.color", config.Default.Color)
	viper.SetDefault("sound.enabled", config.Sound.Enabled)
	viper.SetDefault("sound.file", config.Sound.File)
	viper.SetDefault("sound.volume", config.Sound.Volume)
	viper.SetDefault("display.show_milliseconds", config.Display.ShowMilliseconds)
	viper.SetDefault("display.animation", config.Display.Animation)
	viper.SetDefault("display.color_scheme", config.Display.ColorSheme)
	viper.SetDefault("display.colorl", config.Display.Color)
	viper.SetDefault("notifications.desktop", config.Notifications.Desktop)
	viper.SetDefault("notifications.terminal", config.Notifications.Terminal)

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}

func getDefaultSoundFile() string {
	possiblePaths := []string {
		"/usr/share/sounds/alsa/Front_Left.wav",
		"/System/Library/Sounds/Ping.aiff",
		"C:\\Windows\\Media\\chimes.wav",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func CreateDefaultConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	configPath := filepath.Join(home, ".timer.yaml")

	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists at %s", configPath)
	}

	configContent := `# Go Simple Timer Config

		# Default timer settings
		default:
			sound: true
			desktop_notifications: true
			color: true

		# Sound settings
		sound:
			enabled: true
			file: "" # Leave empty to use system default
			volume: 80

		# Display settings
		display:
			show_milliseconds: false
			animation: true
			color_scheme: "default"
			color: true

		# Notification settings
		notifications:
			desktop: true
			terminal: true
	`

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}

	fmt.Printf("Default configuration created at: %s\n", configPath)
	return nil
}