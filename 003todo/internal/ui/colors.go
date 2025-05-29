package ui

import "github.com/fatih/color"

var ( // color functions
	Red			= color.New(color.FgRed).SprintFunc()
	Green		= color.New(color.FgGreen).SprintFunc()
	Yellow		= color.New(color.FgYellow).SprintFunc()
	Blue		= color.New(color.FgBlue).SprintFunc()
	Magenta		= color.New(color.FgMagenta).SprintFunc()
	Cyan		= color.New(color.FgCyan).SprintFunc()
	White		= color.New(color.FgWhite).SprintFunc()

	// style functions
	Bold		= color.New(color.Bold).SprintFunc()
	Underline	= color.New(color.Underline).SprintFunc()
	Dim			= color.New(color.Faint).SprintFunc()
)

func StatusColor(completed bool) func(a ...interface{}) string {
	if completed {
		return Green
	}
	return Yellow
}

func PriorityColor(priority string) func(a ...interface{}) string {
	switch priority {
	case "High":
		return Red
	case "Medium":
		return Yellow
	case "Low":
		return Green
	default:
		return White
	}
}