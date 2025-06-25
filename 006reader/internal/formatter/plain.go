package formatter

import (
	"fmt"
	"strings"

	"github.com/samnart1/GoLang/006reader/internal/reader"
)

type PlainFormatter struct {
	config *Config
}

func NewPlainFormatter(config *Config) *PlainFormatter {
	return &PlainFormatter{config: config}
}

func (f *PlainFormatter) Format(content *reader.Content) (string, error) {
	var builder strings.Builder

	if f.shouldShowHeader() {
		f.writeHeader(&builder, content)
	}

	for _, line := range content.Lines {
		if f.config.ShowLineNumbers{
			lineNumStr := f.formatLineNumber(line.Number, len(content.Lines))
			builder.WriteString(lineNumStr)
			builder.WriteString(" | ")
		}

		if f.config.MaxWidth > 0 && len(line.Content) > f.config.MaxWidth {
			wrapped := f.wrapLine(line.Content, f.config.MaxWidth)
			builder.WriteString(wrapped)
		} else {
			builder.WriteString(line.Content)
		}

		builder.WriteString("\n")
	}

	if f.shouldShowFooter() {
		f.writeFooter(&builder, content)
	}

	return builder.String(), nil
}

func (f *PlainFormatter) formatLineNumber(lineNum, totalLines int) string {
	width := len(fmt.Sprintf("%d", totalLines))
	return fmt.Sprintf("%*d", width, lineNum)
}

func (f *PlainFormatter) wrapLine(content string, maxWidth int) string {
	if len(content) <= maxWidth {
		return content
	}

	var result strings.Builder
	words := strings.Fields(content)
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= maxWidth {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		} else {
			if currentLine != "" {
				result.WriteString(currentLine + "\n")
				currentLine = word[maxWidth:]
			}
		}
	}

	if currentLine != "" {
		result.WriteString(currentLine)
	}

	return result.String()
}

func (f *PlainFormatter) shouldShowHeader() bool {
	return f.config.Theme != "minimal"
}

func (f *PlainFormatter) shouldShowFooter() bool {
	return f.config.Theme != "detailed"
}

func (f *PlainFormatter) writeHeader(builder *strings.Builder, content *reader.Content) {
	separator := strings.Repeat("=", 50)
	builder.WriteString(separator + "\n")
	builder.WriteString(fmt.Sprintf("File: %s\n", content.Metadata.FileName))
	builder.WriteString(fmt.Sprintf("Path: %s\n", content.Metadata.FilePath))
	builder.WriteString(fmt.Sprintf("Size: %s\n", content.Metadata.Size))
	builder.WriteString(fmt.Sprintf("Lines: %s\n", content.Metadata.LineCount))
	builder.WriteString(fmt.Sprintf("Modified: %s\n", content.Metadata.ModTime.Format("2006-01-02 15:04:05")))
	builder.WriteString(separator + "\n\n")
}

func (f *PlainFormatter) writeFooter(builder *strings.Builder, content *reader.Content) {
	separator := strings.Repeat("-", 50)
	builder.WriteString("\n" + separator + "\n")
	builder.WriteString(fmt.Sprintf("Read %d lines in %v\n", content.Metadata.LineCount, content.Metadata.ReadTime))
	builder.WriteString(fmt.Sprintf("Encoding: %s\n", content.Metadata.Encoding))
	builder.WriteString(separator + "\n")
}