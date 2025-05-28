package task

import (
	"strings"

	"github.com/samnart1/GoLang-Projects/003todo/pkg/errors"
)

func ParsePriority(s string) (Priority, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "low", "l", "1":
		return Low, nil
	case "medium", "m", "2":
		return Medium, nil
	case "high", "h", "3":
		return High, nil
	default:
		return Medium, errors.NewValidationError("priority", "must be low, medium, or high")
	}
}

func PriorityColor(p Priority) string {
	switch p {
	case Low:
		return "green"
	case Medium:
		return "yellow"
	case High:
		return "red"
	default:
		return "white"
	}
}