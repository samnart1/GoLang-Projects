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

