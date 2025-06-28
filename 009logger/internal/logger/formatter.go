package logger

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Formatter interface {
	Format(entry Entry) ([]byte, error)
}

type TextFormatter struct {
	TimestampFormat string
}

func (f *TextFormatter) Format(entry Entry) ([]byte, error) {
	timestamp := entry.Timestamp.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level)

	var line string
	if entry.Source != "" {
		line = fmt.Sprintf("[%s] %s [%s] %s", timestamp, level, entry.Source, entry.Message)
	} else {
		line = fmt.Sprintf("[%s] %s %s", timestamp, level, entry.Message)
	}

	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}

	return []byte(line), nil
}

type JSONFormatter struct {
	TimestampFormat string
}

type JSONEntry struct {
	Timestamp 	string					`json:"timestamp"`
	Level		string					`json:"level"`
	Message		string					`json:"message"`
	Source		string					`json:"source,omitempty"`
	Data 		map[string]interface{}	`json:"data,omitempty"`
}

func (f *JSONFormatter) Format(entry Entry) ([]byte, error) {
	jsonEntry := JSONEntry{
		Timestamp: entry.Timestamp.Format(f.TimestampFormat),
		Level: entry.Level,
		Message: entry.Message,
		Source: entry.Source,
		Data: entry.Data,
	}

	data, err := json.Marshal(jsonEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	data = append(data, '\n')

	return data, nil
}