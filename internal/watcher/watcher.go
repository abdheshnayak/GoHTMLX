package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abdheshnayak/gohtmlx/internal/compiler"
	"github.com/abdheshnayak/gohtmlx/internal/config"
	"github.com/abdheshnayak/gohtmlx/pkg/logger"
)

// Watcher handles file watching and automatic rebuilds
type Watcher struct {
	config   *config.Config
	logger   logger.Logger
	compiler *compiler.Compiler
	stopChan chan bool
}

// New creates a new watcher instance
func New(cfg *config.Config, log logger.Logger) *Watcher {
	return &Watcher{
		config:   cfg,
		logger:   log,
		compiler: compiler.New(cfg, log),
		stopChan: make(chan bool),
	}
}

// Start starts watching for file changes
func (w *Watcher) Start() error {
	w.logger.Info("Starting file watcher...", "dir", w.config.SourceDir)

	// Initial build
	if err := w.compiler.Build(); err != nil {
		w.logger.Error("Initial build failed", "error", err)
		return err
	}

	w.logger.Info("File watcher started. Press Ctrl+C to stop.")

	// Simple polling-based file watcher
	go w.watchLoop()

	// Keep the watcher running
	<-w.stopChan
	return nil
}

// Stop stops the file watcher
func (w *Watcher) Stop() error {
	close(w.stopChan)
	return nil
}

// watchLoop implements a simple polling-based file watcher
func (w *Watcher) watchLoop() {
	fileModTimes := make(map[string]time.Time)

	// Initial scan
	w.scanFiles(fileModTimes)

	ticker := time.NewTicker(time.Duration(w.config.Watch.Debounce) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopChan:
			return
		case <-ticker.C:
			w.scanFiles(fileModTimes)
		}
	}
}

// scanFiles scans for file changes
func (w *Watcher) scanFiles(fileModTimes map[string]time.Time) {
	err := filepath.Walk(w.config.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file should be ignored
		if w.shouldIgnoreFile(path) {
			return nil
		}

		// Check if file has valid extension
		ext := strings.ToLower(filepath.Ext(path))
		if _, ok := w.config.Extensions[ext]; !ok {
			return nil
		}

		modTime := info.ModTime()
		lastModTime, exists := fileModTimes[path]

		if !exists {
			// New file
			fileModTimes[path] = modTime
			w.logger.Debug("New file detected", "file", path)
			w.handleFileChange(path)
		} else if modTime.After(lastModTime) {
			// Modified file
			fileModTimes[path] = modTime
			w.logger.Debug("File modified", "file", path)
			w.handleFileChange(path)
		}

		return nil
	})

	if err != nil {
		w.logger.Error("Error scanning files", "error", err)
	}
}

// handleFileChange processes file change events
func (w *Watcher) handleFileChange(filename string) {
	w.logger.Debug("File changed", "file", filename)

	if err := w.compiler.BuildFile(filename); err != nil {
		w.logger.Error("Failed to rebuild file", "file", filename, "error", err)
	} else {
		w.logger.Info("File rebuilt", "file", filename)
	}
}

// shouldIgnoreFile checks if a file should be ignored based on configuration
func (w *Watcher) shouldIgnoreFile(filename string) bool {
	base := filepath.Base(filename)

	for _, pattern := range w.config.Watch.IgnoreFiles {
		if matched, _ := filepath.Match(pattern, base); matched {
			return true
		}
	}

	return false
}
