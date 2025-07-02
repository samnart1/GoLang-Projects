package notifications

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type TerminalNotifier struct {
	enabled bool
}

func NewTerminalNotifier(enabled bool) *TerminalNotifier {
	return &TerminalNotifier{
		enabled: enabled,
	}
}

func (tn *TerminalNotifier) Show(message string) error {
	if !tn.enabled {
		return nil
	}

	border := strings.Repeat("=", len(message)+4)

	color.Red(border)
	color.Red("| %s |", message)
	color.Red(border)
	fmt.Println()

	fmt.Print("\a")

	return nil
}

// func (tn *TerminalNotifier) ShowBanner(title, message string) error {
// 	if !tn.enabled {
// 		return nil
// 	}

// 	maxLen := len(title)
// 	if len(message) > maxLen {
// 		maxLen = len(message)
// 	}

// 	border := strings.Repeat()
// }