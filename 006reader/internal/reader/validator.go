package reader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ValidationError struct {
	FilePath 	string
	Reason		string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.FilePath, e.Reason)
}

type FileValidator struct {
	maxSize		 int64
	allowedExts	 []string
	blockedExts	 []string
	checkContent bool
}

func NewValidator() *FileValidator {
	return &FileValidator{
		maxSize: 100 * 1024 * 1024,	// 100mb
		allowedExts: []string{".txt", ".log", ".md", ".csv", ".json", ".xml", ".yml"},
		blockedExts: []string{".exe", ".bin", ".so", ".dll"},
		checkContent: true,	
	}
}

func (v *FileValidator) SetMaxSize(size int64) *FileValidator {
	v.maxSize = size
	return v
}

func (v *FileValidator) SetAllowedExtensions(exts []string) *FileValidator {
	v.allowedExts = exts
	return v
}

func (v *FileValidator) SetBlockedExtensions(exts []string) *FileValidator {
	v.blockedExts = exts
	return v
}

func (v *FileValidator) EnableContentCheck(enable bool) *FileValidator {
	v.checkContent = enable
	return v
}

func (v *FileValidator) Validate(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &ValidationError{filePath, "file does not exist"}
		}
		return &ValidationError{filePath, fmt.Sprintf("cannot access file: %v", err)}
	}

	if info.IsDir() {
		return &ValidationError{filePath, "path is a directory, not a file"}
	}

	if v.maxSize > 0 && info.Size() > v.maxSize {
		return &ValidationError{filePath, fmt.Sprintf("file too large: %d bytes (max: %d)", info.Size(), v.maxSize)}
	}

	if err := v.validateExtension(filePath); err != nil {
		return err
	}

	if err := v.validatePermissions(filePath); err != nil {
		return err
	}

	if v.checkContent {
		if err := v.validateContent(filePath); err != nil {
			return err
		}
	}

	return nil
}

func (v *FileValidator) validateExtension(filePath string) error {
	ext := strings.ToLower(filepath.Ext(filePath))

	for _, blocked := range v.blockedExts {
		if ext == strings.ToLower(blocked) {
			return &ValidationError{filePath, fmt.Sprintf("file type not allowed: %s", ext)}
		}
	}

	if len(v.allowedExts) > 0 {
		allowed := false
		for _, allowedExt := range v.allowedExts {
			if ext == strings.ToLower(allowedExt) {
				allowed = true
				break
			}
		}
		if !allowed {
			return &ValidationError{filePath, fmt.Sprintf("file type not in allowed list: %s", ext)}
		}
	}

	return nil
}

func (v *FileValidator) validatePermissions(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return &ValidationError{filePath, fmt.Sprintf("cannot open file: %v", err)}
	}
	file.Close()

	return nil
}

func (v *FileValidator) validateContent(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return &ValidationError{filePath, fmt.Sprintf("cannot open for content validation: %v", err)}
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return &ValidationError{filePath, fmt.Sprintf("cannot read file content: %v", err)}
	}

	for i := 0; i < n;  i++ {
		if buffer[i] == 0 {
			return &ValidationError{filePath, "file appears to be binary"}
		}
	}

	return nil
}

func (v *FileValidator) ValidateMultiple(filePaths []string) map[string]error {
	results := make(map[string]error)

	for _, filePath := range filePaths {
		results[filePath] = v.Validate(filePath)
	}

	return results
}

func GetFileType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".txt":
		return "text"
	case ".log":
		return ".log"
	case ".md":
		return "markdown"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	case ".csv":
		return "csv"
	case ".yml", ".yaml":
		return "yaml"
	default:
		return "unknown"
	}
}