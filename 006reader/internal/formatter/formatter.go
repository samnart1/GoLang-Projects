package formatter

import (
	"fmt"

	"github.com/samnart1/GoLang/006reader/internal/reader"
)

type Formatter interface {
	Format(content *reader.Content) (string, error)
	Name() string
}

type Config struct {
	ShowLineNumbers bool
	MaxWidth		int
	Theme			string
	ColorOutput		bool
}

func New(format string, config *Config) (Formatter, error) {
	if config == nil {
		config = &Config{
			MaxWidth: 120,
			Theme: "default",
			ColorOutput: true,
		}
	}

	switch format {
	case "plain", "text":
		return NewPlainFormatter(config), nil
	case "json":
		return NewJSONFormatter(config), nil
	case "table":
		return NewTableFormatter(config), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func GetAvailableFormats() []string {
	return []string{"plain", "json", "table"}
}