package utils

import (
	"os"
	"strings"
	"path/filepath"
	"fmt"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func GetFileExtension(filePath string) string {
	return strings.ToLower(filepath.Ext(filePath))
}

func GetFileName(filePath string) string {
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, ext)
}

func GetAbsolutePath(filePath string) (string, error) {
	return filepath.Abs(filePath)
}

func IsTextFile(filePath string) bool {
	textExtensions := []string{
		".text", ".log", ".md", ".csv", ".json", ".xml",
		".yml", ".yaml", ".ini", ".cfg", ".conf",
	}

	ext := GetFileExtension(filePath)
	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	return false
}

func SanitizeFileName(fileName string) string {
	invalid := []string{"<", ">", ":", "\"", "|", "?", "*", "/", "\\"}

	sanitized := fileName
	for _, char := range invalid {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}

	return sanitized
}

func FormateFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func CreateTempFile(content string, prefix string) (string, error) {
	tempFile, err := os.CreateTemp("", prefix+"*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	if _, err := tempFile.WriteString(content); err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tempFile.Name(), nil
}