package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/samnart1/golang/009l9gger/internal/config"
)

func TestLogger_Log(t *testing.T) {
	// Create temporary log file
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create test config
	cfg := &config.Config{
		Log: config.LogConfig{
			File:            tmpFile.Name(),
			Level:           "info",
			Format:          "text",
			Output:          "file",
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}

	// Create logger
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Test logging
	entry := Entry{
		Message: "Test message",
		Level:   "info",
		Source:  "test",
	}

	if err := logger.Log(entry); err != nil {
		t.Fatalf("Failed to log entry: %v", err)
	}

	// Read log file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, "Test message") {
		t.Errorf("Log content does not contain expected message")
	}
	if !strings.Contains(logContent, "INFO") {
		t.Errorf("Log content does not contain expected level")
	}
	if !strings.Contains(logContent, "test") {
		t.Errorf("Log content does not contain expected source")
	}
}

func TestLogger_LogLevels(t *testing.T) {
	// Create temporary log file
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create test config with warn level
	cfg := &config.Config{
		Log: config.LogConfig{
			File:            tmpFile.Name(),
			Level:           "warn",
			Format:          "text",
			Output:          "file",
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}

	// Create logger
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Test different log levels
	logger.Debug("Debug message")  // Should not be logged
	logger.Info("Info message")    // Should not be logged
	logger.Warn("Warning message") // Should be logged
	logger.Error("Error message")  // Should be logged

	// Read log file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)
	
	// Check that only warn and error messages are logged
	if strings.Contains(logContent, "Debug message") {
		t.Errorf("Debug message should not be logged")
	}
	if strings.Contains(logContent, "Info message") {
		t.Errorf("Info message should not be logged")
	}
	if !strings.Contains(logContent, "Warning message") {
		t.Errorf("Warning message should be logged")
	}
	if !strings.Contains(logContent, "Error message") {
		t.Errorf("Error message should be logged")
	}
}