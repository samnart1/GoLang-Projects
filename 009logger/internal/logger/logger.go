package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/samnart1/golang/009l9gger/internal/config"
)

type Logger struct {
	config 		*config.Config
	writer		*Writer
	formatter	Formatter
	mu			sync.Mutex
}

type Entry struct {
	Message 	string
	Level		string
	Source		string
	Timestamp 	time.Time
	Data		map[string]interface{}
}

func New(cfg *config.Config) (*Logger, error) {
	writer, err := NewWriter(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create writer: %w", err)
	}

	var formatter Formatter
	switch cfg.Log.Format {
	case "json":
		formatter = &JSONFormatter{
			TimestampFormat: cfg.Log.TimestampFormat,
		}

	default:
		formatter = &TextFormatter{
			TimestampFormat: cfg.Log.TimestampFormat,
		}
	}

	return &Logger{
		config: cfg,
		writer: writer,
		formatter: formatter,
	}, nil
}

func (l *Logger) Log(entry Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	if !l.shouldLog(entry.Level) {
		return nil
	}

	formatted, err := l.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("failed to format log entry: %w", err)
	}

	if _, err := l.writer.Write(formatted); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	return nil
}

func (l *Logger) Debug(message string) error {
	return l.Log(Entry{Message: message, Level: "debug"})
}

func (l *Logger) Info(message string) error {
	return l.Log(Entry{Message: message, Level: "info"})
}

func (l *Logger) Warn(message string) error {
	return l.Log(Entry{Message: message, Level: "warn"})
}

func (l *Logger) Error(message string) error {
	return l.Log(Entry{Message: message, Level: "error"})
}

func (l *Logger) Fatal(message string) error {
	err := l.Log(Entry{Message: message, Level: "fatal"})
	os.Exit(1)
	return err
}

func (l *Logger) Close() error {
	if l.writer != nil {
		return l.writer.Close()
	}

	return nil
}

func (l *Logger) shouldLog(level string) bool {
	levelPriority := map[string]int{
		"debug": 0,
		"info":	1,
		"warn":	2,
		"error": 3,
		"fatal": 4,
	}

	entryPriority, exists := levelPriority[level]
	if !exists {
		return true
	}

	configPriority, exists := levelPriority[l.config.Log.Level]
	if !exists {
		return true
	}

	return entryPriority >= configPriority
}

func (l *Logger) GetWriter() io.Writer {
	return l.writer
}