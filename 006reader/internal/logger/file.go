package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FileLogger struct {
	*Logger
	logFile *os.File
}

func NewFileLogger(logDir, appName string) (*FileLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("%s-%s.log", appName, timestamp)
	logFilePath := filepath.Join(logDir, logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	//create core that writes to file
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(logFile),
		zapcore.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller())

	return &FileLogger{
		Logger: &Logger{Logger: logger},
		logFile: logFile,
	}, nil
}

func (fl *FileLogger) Close() error {
	fl.Logger.Sync()
	return fl.logFile.Close()
}

func (fl *FileLogger) RotateLog() error {
	if err := fl.Logger.Sync(); err != nil {
		return fmt.Errorf("failed to sync logger: %w", err)
	}
	if err := fl.logFile.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %w", err)
	}

	
	dir := filepath.Dir(fl.logFile.Name())
	base := filepath.Base(fl.logFile.Name())
	nameWithoutExt := strings.TrimSuffix(base, filepath.Ext(base))
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	newPath := filepath.Join(dir, fmt.Sprintf("%s-%s.log", nameWithoutExt, timestamp))

	
	if err := os.Rename(fl.logFile.Name(), newPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	
	logFile, err := os.OpenFile(fl.logFile.Name(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(logFile),
		zapcore.InfoLevel,
	)

	
	fl.Logger = &Logger{Logger: zap.New(core, zap.AddCaller())}
	fl.logFile = logFile

	return nil
}