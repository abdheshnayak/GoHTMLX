package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	WithFields(fields map[string]interface{}) Logger
}

// slogLogger implements Logger using slog
type slogLogger struct {
	logger *slog.Logger
	fields map[string]interface{}
}

// New creates a new logger instance
func New() Logger {
	// Use charmbracelet/log for better formatting
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "GoHTMLX",
	})

	return &charmLogger{logger: logger}
}

// charmLogger implements Logger using charmbracelet/log
type charmLogger struct {
	logger *log.Logger
	fields map[string]interface{}
}

func (l *charmLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *charmLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *charmLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *charmLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

func (l *charmLogger) WithFields(fields map[string]interface{}) Logger {
	newLogger := &charmLogger{
		logger: l.logger,
		fields: make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// SetLevel sets the logging level
func SetLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}
}
