package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/samnart1/golang/009l9gger/internal/config"
)

type Writer struct {
	fileWriter 		io.WriteCloser
	consoleWriter	io.Writer
	output			string
}

func NewWriter(cfg *config.Config) (*Writer, error) {
	w := &Writer{
		output: cfg.Log.Output,
		consoleWriter: os.Stdout,
	}

	if cfg.Log.Output == "file" || cfg.Log.Output == "both" {
		if cfg.Log.File == "" {
			return nil, fmt.Errorf("log file path is required for file output")
		}

		file, err := os.OpenFile(cfg.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		w.fileWriter = file
	}

	return w, nil
}

func (w *Writer) Write(data []byte) (int, error) {
	var err error
	var totalBytes int

	if w.output == "file" || w.output == "both" {
		if w.fileWriter != nil {
			n, writeErr := w.fileWriter.Write(data)
			totalBytes += n
			if writeErr != nil {
				err = fmt.Errorf("failed to write to file: %w", writeErr)
			}
		}
	}

	if w.output == "console" || w.output == "both" {
		n, writeErr := w.consoleWriter.Write(data)
		totalBytes += n
		if writeErr != nil {
			if err != nil {
				err = fmt.Errorf("%v; failed to write to console: %w", err, writeErr)
			} else {
				err = fmt.Errorf("failed to write to console: %w", writeErr)
			}
		}
	}

	return totalBytes, err
}

func (w *Writer) Close() error {
	if w.fileWriter != nil {
		return w.fileWriter.Close()
	}
	return nil
}