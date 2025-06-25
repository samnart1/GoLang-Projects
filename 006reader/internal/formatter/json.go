package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/samnart1/GoLang/006reader/internal/reader"
)

type JSONFormatter struct {
	config *Config
}

func NewJSONFormatter(config *Config) *JSONFormatter {
	return &JSONFormatter{config: config}
}

func (f *JSONFormatter) Format(content *reader.Content) (string, error) {
	output := struct {
		*reader.Content
		FormatterInfo map[string]interface{} `json:"formatter_info"`
	}{
		Content: content,
		FormatterInfo: map[string]interface{}{
			"format":           "json",
			"show_line_numbers": f.config.ShowLineNumbers,
			"max_width":        f.config.MaxWidth,
			"theme":            f.config.Theme,
		},
	}

	jsonBytes, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

func (f *JSONFormatter) Name() string {
	return "json"
}