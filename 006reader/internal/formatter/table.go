package formatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samnart1/GoLang/006reader/internal/reader"
)

type TableFormatter struct {
	config	*Config
}

func NewTableFormatter(config *Config) *TableFormatter {
	return &TableFormatter{config: config}
}

func (f *TableFormatter) Format(content *reader.Content) (string, error) {
	if len(content.Lines) == 0 {
		return "No content to display\n", nil
	}

	var builder strings.Builder

	lineNumWidth := len(strconv.Itoa(len(content.Lines)))
	contentWidth := f.calculateContentWidth(content.Lines)

	if f.config.MaxWidth > 0 {
		availableWidth := f.config.MaxWidth - lineNumWidth - 7
		if availableWidth > 0 && contentWidth > availableWidth {
			contentWidth = availableWidth
		}
	}

	f.writeTableHeader(&builder, lineNumWidth, contentWidth)

	for _, line := range content.Lines {
		f.writeTableRow(&builder, line, lineNumWidth, contentWidth)
	}

	f.writeTableFooter(&builder, lineNumWidth, contentWidth, content)

	return builder.String(), nil
}

func (f *TableFormatter) calculateContentWidth(lines []reader.Line) int {
	maxWidth := 20

	for _, line := range lines {
		if len(line.Content) > maxWidth {
			maxWidth = len(line.Content)
		}
	}

	if maxWidth > 100 {
		maxWidth = 100
	}

	return maxWidth
}

func (f *TableFormatter) writeTableHeader(builder *strings.Builder, lineWidth, contentWidth int) {
	//top border
	builder.WriteString("r")
	builder.WriteString(strings.Repeat("-", lineWidth+2))
	builder.WriteString("T")
	builder.WriteString(strings.Repeat("-", contentWidth+2))
	builder.WriteString("┐\n")

	//header row
	builder.WriteString("|")
	builder.WriteString(fmt.Sprintf(" %-*s ", lineWidth, "Line"))
	builder.WriteString("|")
	builder.WriteString(fmt.Sprintf(" %-*s ", contentWidth, "Content"))
	builder.WriteString("|\n")

	//header separator
	builder.WriteString("├")
	builder.WriteString(strings.Repeat("─", lineWidth+2))
	builder.WriteString("┼")
	builder.WriteString(strings.Repeat("─", contentWidth+2))
	builder.WriteString("┤\n")
}

func (f *TableFormatter) writeTableRow(builder *strings.Builder, line reader.Line, lineWidth, contentWidth int) {
	content := line.Content

	if len(content) > contentWidth {
		content = content[:contentWidth-3] + "..."
	}

	builder.WriteString("|")
	builder.WriteString(fmt.Sprintf(" %*d ", lineWidth, line.Number))
	builder.WriteString("|")
	builder.WriteString(fmt.Sprintf(" %-*s ", contentWidth, content))
	builder.WriteString("|\n")
}

func (f *TableFormatter) writeTableFooter(builder *strings.Builder, lineWidth, contentWidth int, content *reader.Content) {
	//bottom border
	builder.WriteString("└")
	builder.WriteString(strings.Repeat("─", lineWidth+2))
	builder.WriteString("┴")
	builder.WriteString(strings.Repeat("─", contentWidth+2))
	builder.WriteString("┘\n")

	builder.WriteString(fmt.Sprintf("\nSummary: %d lines, %d bytes, read in %v\n", 
		content.Metadata.LineCount,
		content.Metadata.Size,
		content.Metadata.ReadTime))
}