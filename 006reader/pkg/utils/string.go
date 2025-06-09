package utils

import "strings"

func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	if maxLength <= 3 {
		return s[:maxLength]
	}

	return s[:maxLength-3] + "..."
}

func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func CountLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func RemoveEmptyLines(s string) string {
	lines := strings.Split(s, "\n")
	var nonEmptyLines []string

	for _, line := range lines {
		if !IsEmpty(line) {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	return strings.Join(nonEmptyLines, "\n")
}

func IndentString(s, indent string) string {
	lines := strings.Split(s, "\n")
	var indentedLines []string

	for _, line := range lines {
		indentedLines = append(indentedLines, indent+line)
	}

	return strings.Join(indentedLines, "\n")
}