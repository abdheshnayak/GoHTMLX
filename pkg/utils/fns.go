package utils

import (
	"log/slog"
	"strings"
	"text/template"
)

// Logger is used for progress and errors. Default is no-op; set Log to inject (e.g. from CLI or integration).
type Logger interface {
	Info(msg string, kvs ...any)
	Error(msg string, kvs ...any)
}

type noopLogger struct{}

func (noopLogger) Info(msg string, kvs ...any)  {}
func (noopLogger) Error(msg string, kvs ...any) {}

// Log is the global logger. Default is no-op; set to a real logger (e.g. utils.NewSlogLogger(slog.Default())) from main or example.
var Log Logger = noopLogger{}

// NewSlogLogger returns a Logger that forwards to the given slog.Logger.
func NewSlogLogger(l *slog.Logger) Logger {
	return &slogLogger{l: l}
}

type slogLogger struct{ l *slog.Logger }

func (s *slogLogger) Info(msg string, kvs ...any)  { s.l.Info(msg, kvs...) }
func (s *slogLogger) Error(msg string, kvs ...any) { s.l.Error(msg, kvs...) }

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func GetSections(tmpl *template.Template) []string {
	sectionNames := []string{}
	t := tmpl.Templates()
	for _, v := range t {
		sectionNames = append(sectionNames, v.Name())
	}

	return sectionNames
}

func ParseSections(tmpl *template.Template) (map[string]string, error) {
	sectionNames := GetSections(tmpl)
	t := tmpl.Templates()
	for _, v := range t {
		sectionNames = append(sectionNames, v.Name())
	}

	sections := make(map[string]string)

	for _, section := range sectionNames {
		var buffer strings.Builder
		err := tmpl.ExecuteTemplate(&buffer, section, nil)
		if err != nil {
			return nil, err
		}
		sections[section] = buffer.String()
	}

	return sections, nil
}
