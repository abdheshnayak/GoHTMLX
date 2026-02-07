// Package validate checks GoHTMLX .html files for unclosed or mismatched
// comment blocks (<!-- + define --> / <!-- + end -->, etc.).
package validate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reOpen  = regexp.MustCompile(`<!--\s*([*+|])\s+define\s+"([^"]*)"\s*-->`)
	reClose = regexp.MustCompile(`<!--\s*([*+|])\s+end\s*-->`)
)

// Run walks src for .html files and reports comment-structure errors to w.
// Returns nil if no errors, or an error summarizing the first failure.
func Run(src string, w io.Writer) error {
	var hadErr bool
	var firstErr string
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			hadErr = true
			if firstErr == "" {
				firstErr = fmt.Sprintf("%s: read: %v", path, err)
			}
			fmt.Fprintf(w, "%s: read: %v\n", path, err)
			return nil
		}
		lines := strings.Split(string(content), "\n")
		var stack []string
		for i, line := range lines {
			lineNum := i + 1
			for _, m := range reClose.FindAllStringSubmatch(line, -1) {
				if len(m) < 2 {
					continue
				}
				delim := m[1]
				if len(stack) == 0 {
					msg := fmt.Sprintf("%s:%d: unexpected <!-- %s end --> (no open block)\n", path, lineNum, delim)
					if firstErr == "" {
						firstErr = strings.TrimSpace(msg)
					}
					fmt.Fprint(w, msg)
					hadErr = true
					continue
				}
				top := stack[len(stack)-1]
				if top != delim {
					msg := fmt.Sprintf("%s:%d: <!-- %s end --> does not match open <!-- %s define -->\n", path, lineNum, delim, top)
					if firstErr == "" {
						firstErr = strings.TrimSpace(msg)
					}
					fmt.Fprint(w, msg)
					hadErr = true
				}
				stack = stack[:len(stack)-1]
			}
			for _, m := range reOpen.FindAllStringSubmatch(line, -1) {
				if len(m) >= 2 {
					stack = append(stack, m[1])
				}
			}
		}
		if len(stack) > 0 {
			msg := fmt.Sprintf("%s:%d: unclosed block(s) <!-- %s define --> (missing <!-- %s end -->)\n", path, len(lines), strings.Join(stack, ", "), stack[len(stack)-1])
			if firstErr == "" {
				firstErr = strings.TrimSpace(msg)
			}
			fmt.Fprint(w, msg)
			hadErr = true
		}
		return nil
	})
	if err != nil {
		return err
	}
	if hadErr {
		return fmt.Errorf("validation failed: %s", firstErr)
	}
	return nil
}
