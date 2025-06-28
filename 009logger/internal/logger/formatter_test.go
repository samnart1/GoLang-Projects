package logger

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTextFormatter_Format(t *testing.T) {
	formatter := &TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Level:     "info",
		Message:   "Test message",
		Source:    "test",
	}

	result, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format entry: %v", err)
	}

	output := string(result)
	expected := "[2024-01-01 12:00:00] INFO [test] Test message\n"
	
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	formatter := &JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Level:     "info",
		Message:   "Test message",
		Source:    "test",
		Data:      map[string]interface{}{"key": "value"},
	}

	result, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("Failed to format entry: %v", err)
	}

	// Parse JSON to verify structure
	var jsonEntry JSONEntry
	if err := json.Unmarshal(result, &jsonEntry); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if jsonEntry.Level != "info" {
		t.Errorf("Expected level 'info', got %q", jsonEntry.Level)
	}
	if jsonEntry.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got %q", jsonEntry.Message)
	}
	if jsonEntry.Source != "test" {
		t.Errorf("Expected source 'test', got %q", jsonEntry.Source)
	}
}
