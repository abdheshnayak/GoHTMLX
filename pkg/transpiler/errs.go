package transpiler

import (
	"fmt"
	"strings"
)

// TranspileError is returned when transpilation fails. It includes source location
// so the caller can show file:line and optional snippet.
type TranspileError struct {
	Component string // component name, if applicable
	FilePath  string
	Line      int
	Message   string
	Snippet   string
}

func (e *TranspileError) Error() string {
	if e.Line > 0 {
		if e.Snippet != "" {
			return fmt.Sprintf("%s:%d: %s\n%s", e.FilePath, e.Line, e.Message, e.Snippet)
		}
		return fmt.Sprintf("%s:%d: %s", e.FilePath, e.Line, e.Message)
	}
	if e.FilePath != "" {
		if e.Snippet != "" {
			return fmt.Sprintf("%s: %s\n%s", e.FilePath, e.Message, e.Snippet)
		}
		return fmt.Sprintf("%s: %s", e.FilePath, e.Message)
	}
	return e.Message
}

func lineForComponent(content []byte, name string) int {
	search := `define "` + name + `"`
	s := string(content)
	idx := strings.Index(s, search)
	if idx < 0 {
		return 0
	}
	return 1 + strings.Count(s[:idx], "\n")
}

func snippetAtLine(content []byte, line int, contextLines int) string {
	lines := strings.Split(string(content), "\n")
	if line < 1 || line > len(lines) {
		return ""
	}
	start := line - 1 - contextLines
	if start < 0 {
		start = 0
	}
	end := line + contextLines
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[start:end], "\n")
}

func wrapTranspileErr(component, filePath string, fileContent []byte, err error) error {
	if err == nil {
		return nil
	}
	line := 0
	var snippet string
	if len(fileContent) > 0 {
		line = lineForComponent(fileContent, component)
		if line > 0 {
			snippet = snippetAtLine(fileContent, line, 2)
		}
	}
	return &TranspileError{
		Component: component,
		FilePath:  filePath,
		Line:      line,
		Message:   err.Error(),
		Snippet:   snippet,
	}
}
