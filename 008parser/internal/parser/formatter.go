package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/samnart1/golang/008parser/internal/config"
	"github.com/samnart1/golang/008parser/pkg/types"
)

type Formatter struct {
	config *config.Config
}

func NewFormatter(cfg *config.Config) *Formatter {
	return &Formatter{
		config: cfg,
	}
}

func (f *Formatter) Format(data interface{}, options types.OutputOptions) (string, error) {
	switch options.Format {
	case "compact":
		return f.formatCompact(data)
	case "pretty":
		return f.formatPretty(data, options)
	case "tree":
		return f.formatTree(data, options, 0)
	case "table":
		return f.formatTable(data, options)
	default:
		return "", fmt.Errorf("unsupported format: %s", options.Format)
	}
}

func (f *Formatter) formatCompact(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func (f *Formatter) formatPretty(data interface{}, options types.OutputOptions) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if options.UseTabs {
		encoder.SetIndent("", "\t")
	} else {
		indent := strings.Repeat(" ", options.Indent)
		encoder.SetIndent("", indent)
	}

	if options.SortKeys {
		data = f.sortKeys(data)
	}

	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	result := buf.String()
	result = strings.TrimSuffix(result, "\n")

	if options.Colors {
		result = f.colorizeJSON(result)
	}

	return result, nil
}

func (f *Formatter) formatTree(data interface{}, options types.OutputOptions, depth int) (string, error) {
	if options.MaxDepth > 0 && depth >= options.MaxDepth {
		return "...", nil
	}

	var result strings.Builder
	indent := strings.Repeat(" ", depth)

	switch obj := data.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(obj))
		for key := range obj {
			keys = append(keys, key)
		}

		if options.SortKeys {
			sort.Strings(keys)
		}

		for i, key := range keys {
			if i > 0 {
				result.WriteString("\n")
			}

			result.WriteString(indent)

			if options.Colors {
				result.WriteString(color.BlueString("├── %s", key))
			} else {
				result.WriteString(fmt.Sprintf("├── %s", key))
			}

			if options.ShowTypes {
				valueType := getJSONType(obj[key])
				if options.Colors {
					result.WriteString(color.YellowString(" (%s)", valueType))
				} else {
					result.WriteString(fmt.Sprintf(" (%s)", valueType))
				}
			}

			if isPrimitive(obj[key]) {
				result.WriteString(": ")
				valueStr := f.formatPrimitiveValue(obj[key], options.Colors)
				result.WriteString(valueStr)
			} else {
				result.WriteString("\n")
				childTree, err := f.formatTree(obj[key], options, depth+1)
				if err != nil {
					return "", err
				}
				result.WriteString(childTree)
			}
		}

	case []interface{}:
		for i, item := range obj {
			if i > 0 {
				result.WriteString("\n")
			}

			result.WriteString(indent)

			if options.Colors {
				result.WriteString(color.BlueString("├── [%d]", i))
			} else {
				result.WriteString(fmt.Sprintf("├── [%d]", i))
			}

			if options.ShowTypes {
				valueType := getJSONType(item)
				if options.Colors {
					result.WriteString(color.YellowString(" (%s)", valueType))
				} else {
					result.WriteString(fmt.Sprintf(" (%s)", valueType))
				}
			}

			if isPrimitive(item) {
				result.WriteString(": ")
				valueStr := f.formatPrimitiveValue(item, options.Colors)
				result.WriteString(valueStr)
			} else {
				result.WriteString("\n")
				childTree, err := f.formatTree(item, options, depth+1)
				if err != nil {
					return "", err
				}
				result.WriteString(childTree)
			}
		}

	default:
		valueStr := f.formatPrimitiveValue(data, options.Colors)
		result.WriteString(indent + valueStr)		
	}

	return result.String(), nil
}

func (f *Formatter) formatTable(data interface{}, options types.OutputOptions) (string, error) {
	arr, ok := data.([]interface{})
	if !ok {
		return "", fmt.Errorf("table format only supports arrays")
	}

	if len(arr) == 0 {
		return "Empty array", nil
	}

	for _, item := range arr {
		if _, ok := item.(map[string]interface{}); !ok {
			return "", fmt.Errorf("table format requires array of objects")
		}
	}

	keySet := make(map[string]bool)
	for _, item := range arr {
		obj := item.(map[string]interface{})
		for key := range obj {
			keySet[key] = true
		}
	}

	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	colWidths := make(map[string]int)
	for _, key := range keys {
		colWidths[key] = len(key)
	}

	for _, item := range arr {
		obj := item.(map[string]interface{})
		for _, key := range keys {
			if value, exists := obj[key]; exists {
				valueStr := fmt.Sprintf("%v", value)
				if len(valueStr) > colWidths[key] {
					colWidths[key] = len(valueStr)
				}
			}
		}
	}

	var result strings.Builder

	for i, key := range keys {
		if i > 0 {
			result.WriteString(" | ")
		}
		if options.Colors {
			result.WriteString(color.CyanString("%-*s", colWidths[key], key))
		} else {
			result.WriteString(fmt.Sprintf("%-*s", colWidths[key], key))
		}
	}
	result.WriteString("\n")

	for i, key := range keys {
		if i > 0 {
			result.WriteString("-+-")
		}
		result.WriteString(strings.Repeat("-", colWidths[key]))
	}
	result.WriteString("\n")

	for _, item := range arr {
		obj := item.(map[string]interface{})
		for i, key := range keys {
			if i > 0 {
				result.WriteString(" | ")
			}

			var valueStr string
			if value, exists := obj[key]; exists {
				valueStr = fmt.Sprintf("%v", value)
			} else {
				valueStr = ""
			}

			if options.Colors {
				result.WriteString(color.WhiteString("%-*s", colWidths[key], valueStr))
			} else {
				result.WriteString(fmt.Sprintf("%-*s", colWidths[key], valueStr))
			}
		}
		result.WriteString("\n")
	}
	return result.String(), nil
}

//helper functions
func (f *Formatter) sortKeys(data interface{}) interface{} {
	switch obj := data.(type) {
	case map[string]interface{}:
		sorted := make(map[string]interface{})
		keys := make([]string, 0, len(obj))
		for key := range obj {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			sorted[key] = f.sortKeys(obj[key])
		}
		return sorted

	case []interface{}:
		sorted := make([]interface{}, len(obj))
		for i, item := range obj {
			sorted[i] = f.sortKeys(item)
		}
		return sorted

	default:
		return data
	}
}

func (f *Formatter) colorizeJSON(jsonStr string) string {
	jsonStr = color.GreenString(jsonStr)

	return jsonStr
}

func (f *Formatter) formatPrimitiveValue(value interface{}, useColors bool) string {
	switch v := value.(type) {
	case string:
		if useColors {
			return color.GreenString(`"%s"`, v)
		}
		return fmt.Sprintf(`"%s"`, v)
		
	case float64:
		if useColors {
			return color.MagentaString("%.6g", v)
		}
		return fmt.Sprintf("%.6g", v)

	case bool:
		if useColors {
			return color.YellowString("%t", v)
		}
		return fmt.Sprintf("%t", v)

	case nil:
		if useColors {
			return color.RedString("null")
		}
		return "null"

	default:
		return fmt.Sprintf("%v", v)
	}
}

func isPrimitive(value interface{}) bool {
	switch value.(type) {
	case string, float64, bool, nil: 
		return true

	default:
		return false
	}
}

// func getJSONType(data interface{}) string {
// 	switch data.(type) {
// 	case map[string]interface{}:
// 		return "object"

// 	case []interface{}:
// 		return "array"

// 	case string:
// 		return "string"

// 	case float64:
// 		return "number"

// 	case bool:
// 		return "boolean"

// 	case nil:
// 		return "null"

// 	default:
// 		return "unknown"
// 	}
// }