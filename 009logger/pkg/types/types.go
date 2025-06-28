package types

import "time"

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo LogLevel = "info"
	LogLevelWarn LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormattJSON LogFormat = "json"
)

type LogOutput string

const (
	LogOutputFile LogOutput = "file"
	LogOutputConsole LogOutput = "console"
	LogOutputBoth LogOutput = "both"
)

type LogEntry struct {
	Timestamp 	time.Time				`json:"timestamp"`
	Level 		LogLevel				`json:"level"`
	Message		string					`json:"message"`
	Source		string					`json:"source,omitempty"`
	Data		map[string]interface{}	`json:"data,omitempty"`
}

type APIResponse struct {
	Success		bool		`json:"success"`
	Data		interface{}	`json:"data,omitempty"`
	Error 		string		`json:"error,omitempty"`
	Timestamp	time.Time	`json:"timestamp"`
}

type BatchLogRequest struct {
	Entries []LogEntry `json:"entries"`
}

type LogStats struct {
	TotalEntries	int64				`json:"total_entries"`
	EntriesByLevel	map[LogLevel]int64	`json:"entries_by_level"`
	LastEntry		time.Time			`json:"last_entry"`
	StartTime		time.Time			`json:"start_time"`
}