package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samnart1/golang/008parser/internal/config"
	"github.com/samnart1/golang/008parser/pkg/types"
)

type Parser struct {
	config *config.Config
}

func New(cfg *config.Config) *Parser {
	return &Parser{
		config: cfg,
	}
}

func (p *Parser) ParseFile(filePath string) (*types.ParseResult, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot access file: %w", err)
	}

	if fileInfo.Size() > p.config.MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (limit: %d bytes)",
			fileInfo.Size(), p.config.MaxFileSize)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	return p.ParseBytes(data, filePath)
}

func (p *Parser) ParseBytes(data []byte, source string) (*types.ParseResult, error) {
	var jsonData interface{}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	result := &types.ParseResult{
		Data: jsonData,
		FilePath: source,
		Size: int64(len(data)),
		Type: getJSONType(jsonData),
	}

	if obj, ok := jsonData.(map[string]interface{}); ok {
		result.KeyCount = len(obj)
	} else if arr, ok := jsonData.([]interface{}); ok {
		result.ArrayLength = len(arr)
	}

	return result, nil
}

func (p *Parser) ParseString(jsonStr, source string) (*types.ParseResult, error) {
	return p.ParseBytes([]byte(jsonStr), source)
}

func getJSONType(data interface{}) string {
	switch data.(type) {
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	default:
		return "unknown"
	}
}