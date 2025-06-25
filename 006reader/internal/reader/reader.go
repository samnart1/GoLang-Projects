package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Content struct {
	Lines		[]Line		`json:"lines"`
	Metadata	Metadata	`json:"metadata"`
}

type Line struct {
	Number		int		`json:"number"`
	Content		string	`json:"content"`
}

type Metadata struct {
	FilePath	string			`json:"file_path"`
	FileName	string			`json:"file_name"`
	Size		int64			`json:"size"`
	ModTime		time.Time		`json:"mod_time"`
	LineCount	int				`json:"line_count"`
	Encoding	string			`json:"encoding"`
	ReadTime	time.Duration	`json:"read_time"`
}

type Config struct {
	FilePath	string
	MaxLines	int
	ShowLines	bool
	Encoding	string
	BufferSize	int
}

type Reader struct {
	config *Config
}

func New(config *Config) *Reader {
	if config.BufferSize <= 0 {
		config.BufferSize = 8192
	}

	return &Reader{
		config: config,
	}
}

func (r *Reader) Read() (*Content, error) {
	startTime := time.Now()

	file, err := os.Open(r.config.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	lines, err := r.readLines(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read lines: %w", err)
	}

	content := &Content{
		Lines: lines,
		Metadata: Metadata{
			FilePath: r.config.FilePath,
			FileName: filepath.Base(r.config.FilePath),
			Size: info.Size(),
			ModTime: info.ModTime(),
			LineCount: len(lines),
			Encoding: r.config.FilePath,
			ReadTime: time.Since(startTime),
		},
	}

	return content, nil
}

func (r *Reader) readLines(file *os.File) ([]Line, error) {
	var lines []Line
	scanner := bufio.NewScanner(file)

	buf := make([]byte, 0, r.config.BufferSize)
	scanner.Buffer(buf, r.config.BufferSize*2)

	lineNum := 1
	for scanner.Scan() {
		if r.config.MaxLines > 0 && lineNum > r.config.MaxLines {
			break
		}

		line := Line{
			Number: lineNum,
			Content: scanner.Text(),
		}
		lines = append(lines, line)
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

func (r *Reader) ReadStream(callback func(line Line) error) error {
	file, err := os.Open(r.config.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, r.config.BufferSize)
	scanner.Buffer(buf, r.config.BufferSize*2)

	lineNum := 1
	for scanner.Scan() {
		if r.config.MaxLines > 0 && lineNum > r.config.MaxLines {
			break
		}

		line := Line{
			Number: lineNum,
			Content: scanner.Text(),
		}

		if err := callback(line); err != nil {
			return fmt.Errorf("callback error at line %d: %w", lineNum, err)
		}

		lineNum++
	}
	return scanner.Err()
}

func GetFileInfo(filePath string) (*Metadata, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &Metadata{
		FilePath: filePath,
		FileName: filepath.Base(filePath),
		Size: info.Size(),
		ModTime: info.ModTime(),
	}, nil
}

func ValidateFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("cannot access file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a Directory: %s", filePath)
	}

	// let's try and open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	file.Close()

	return nil
}

func DetectEncoding(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	if n >= 3 && buffer[0] == 0xEF && buffer[1] == 0xBB && buffer[2] == 0xBF {
		return "utf-8-bom", nil
	}

	if isValidUTF8(buffer[:n]) {
		return "utf-8", nil
	}

	return "unknown", nil
}

func isValidUTF8(data []byte) bool {
	return strings.ToValidUTF8(string(data), "") == string(data)
}