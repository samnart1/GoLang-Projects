package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/samnart1/golang/008parser/internal/config"
	"github.com/samnart1/golang/008parser/pkg/types"
)

type Validator struct {
	config *config.Config
}

func NewValidator(cfg *config.Config) *Validator {
	return &Validator{
		config: cfg,
	}
}

func (v *Validator) ValidateFile(filePath string, strict bool) types.ValidationResult {
	result := types.ValidationResult{
		FilePath: filePath,
		IsValid: false,
		Errors: []string{},
		Warnings: []string{},
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("cannot access file: %v", err))
		return result
	}

	if fileInfo.Size() == 0 {
		result.Errors = append(result.Errors, "file is empty")
		return result
	}

	if fileInfo.Size() > v.config.MaxFileSize {
		result.Errors = append(result.Errors, fmt.Sprintf("file too large: %d bytes (limit: %d bytes)",
			fileInfo.Size(), v.config.MaxFileSize))
		return result
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("cannot read file: %v", err))
		return result
	}

	return v.ValidateBytes(data, strict)
}

func (v *Validator) ValidateBytes(data []byte, strict bool) types.ValidationResult {
	result := types.ValidationResult{
		IsValid: false,
		Errors: []string{},
		Warnings: []string{},
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("invalid JSON syntax: %v", err))
		return result
	}

	if strict {
		v.performStrictValidation(data, &result)
	}

	v.checkPerformanceWarnings(jsonData, &result)

	if len(result.Errors) == 0 {
		result.IsValid = true
	}

	return result
}

func (v *Validator) ValidateString(jsonStr string, strict bool) types.ValidationResult {
	return v.ValidateBytes([]byte(jsonStr), strict)
}

func (v *Validator) performStrictValidation(data []byte, result *types.ValidationResult) {
	dataStr := string(data)

	if v.hasDuplicateKeys(dataStr) {
		result.Errors = append(result.Errors, "duplicate keys detected")
	}

	if v.hasTrailingCommas(dataStr) {
		result.Warnings = append(result.Warnings, "trailing commas detected (not standard JSON)")
	}

	if v.hasComments(dataStr) {
		result.Warnings = append(result.Warnings, "comments detected (not allowed in strict JSON)")
	}
}

func (v *Validator) checkPerformanceWarnings(data interface{}, result *types.ValidationResult) {
	depth := v.calculateDepth(data, 0)
	if depth > 20 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("deeply nested strucutre (depth: %d)", depth))
	}

	size := v.calculateSize(data)
	if size > 1000 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("large object count: %d", size))
	}
}

func (v *Validator) hasDuplicateKeys(jsonStr string) bool {
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.DisallowUnknownFields()

	var result interface{}
	err := decoder.Decode(&result)

	return err != nil && strings.Contains(err.Error(), "duplicate")
}

func (v *Validator) hasTrailingCommas(jsonStr string) bool {
	return strings.Contains(jsonStr, ",}") || strings.Contains(jsonStr, ",]")
}

func (v *Validator) hasComments(jsonStr string) bool {
	return strings.Contains(jsonStr, "//") || strings.Contains(jsonStr, "/*")
}

func (v *Validator) calculateDepth(data interface{}, currentDepth int) int {
	maxDepth := currentDepth

	switch obj := data.(type) {
	case map[string]interface{}:
		for _, value := range obj {
			depth := v.calculateDepth(value, currentDepth+1)
			if depth > maxDepth {
				maxDepth = depth
			}
		}

	case []interface{}:
		for _, value := range obj {
			depth := v.calculateDepth(value, currentDepth+1)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}

	return maxDepth
}

func (v *Validator) calculateSize(data interface{}) int {
	switch obj := data.(type) {
	case map[string]interface{}:
		size := len(obj)
		for _, value := range obj {
			size += v.calculateSize(value)
		}
		return size

	case []interface{}:
		size := len(obj)
		for _, value := range obj {
			size += v.calculateSize(value)
		}
		return size

	default:
		return 1
	}
}