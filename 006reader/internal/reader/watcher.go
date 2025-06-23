package reader

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WatchEvent struct {
	Type 	string		`json:"type"`
	Content	string		`json:"content"`
	Time	time.Time	`json:"time"`
	LineNum	int			`json:"line_num,omitempty"`
}

type WatchConfig struct {
	FilePath		string
	TailMode		bool
	PollInterval	int
	BufferSize		int
}

type Watcher struct {
	config		*WatchConfig
	watcher		*fsnotify.Watcher
	events		chan WatchEvent
	errors 		chan error
	lastSize 	int64
	lastLineNum	int
	done		chan bool
}

func NewWatcher(config *WatchConfig) *Watcher {
	return &Watcher{
		config: config,
		events: make(chan WatchEvent, 100),
		errors: make(chan error, 10),
		done: 	make(chan bool),
	}
}

func (w *Watcher) Start() error {
	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	if err := w.watcher.Add(w.config.FilePath); err != nil {
		return fmt.Errorf("failed to watch file: %w", err)
	}

	if info, err := os.Stat(w.config.FilePath); err == nil {
		w.lastSize = info.Size()
		if w.config.TailMode {
			// in tail mode, lets count existing lines and start from the bottom
			w.lastLineNum = w.countLines()
		}
	}

	go w.processEvents()

	if !w.config.TailMode {
		w.sendInitialContent()
	}

	return nil
}

func (w *Watcher) Stop() {
	if w.watcher != nil {
		w.watcher.Close()
	}
	close(w.done)
}

func (w *Watcher) Events() <-chan WatchEvent {
	return w.events
}

func (w *Watcher) Errors() <-chan error {
	return w.errors
}

func (w *Watcher) processEvents() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleFileEvent(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			w.errors <- err

		case <-w.done:
			return
		}
	}
}

func (w *Watcher) handleFileEvent(event fsnotify.Event) {
	if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
		w.handleFileWrite()
	} else if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
		w.events <- WatchEvent{
			Type: "removed",
			Content: "File was removed or renamed",
			Time: time.Now(),
		}
	}
}

func (w *Watcher) handleFileWrite() {
	info, err := os.Stat(w.config.FilePath)
	if err != nil {
		w.errors <- fmt.Errorf("failed to stat file: %w", err)
		return
	}

	currentSize := info.Size()
	if currentSize == w.lastSize {
		return
	}

	//read new content
	if w.config.TailMode {
		w.readNewLines()
	} else {
		w.readFullFile()
	}

	w.lastSize = currentSize
}

func (w *Watcher) readNewLines() {
	reader := New(&Config{
		FilePath: w.config.FilePath,
		BufferSize: w.config.BufferSize,
	})

	content, err := reader.Read()
	if err != nil {
		w.errors <- fmt.Errorf("failed to read file: %w", err)
		return
	}

	for i := w.lastLineNum; i < len(content.Lines); i++ {
		line := content.Lines[i]
		w.events <- WatchEvent{
			Type: "new_line",
			Content: line.Content,
			Time: time.Now(),
			LineNum: line.Number,
		}
	}
}