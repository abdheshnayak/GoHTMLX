package transpiler

import (
	"fmt"
	"strings"
)

// TranspileError is returned when transpilation fails. Use it so the caller can show
// file:line and an optional snippet. Use errors.As(err, &te) to detect and inspect it.
//
//   - FilePath: source .html file or generated .go file path, when available.
//   - Line: approximate line in that file (e.g. component start); 0 if unknown.
//   - Message: short error description.
//   - Snippet: optional context lines around the error.
//   - Component: component name when the error is tied to a specific component.
type TranspileError struct {
	Component string // component name, if applicable
	FilePath  string // source or generated file path
	Line      int    // line number, or 0 if unknown
	Message   string
	Snippet   string // optional context lines
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
