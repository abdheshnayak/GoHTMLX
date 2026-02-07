// validate is a small checker for GoHTMLX .html component files. It reports
// unclosed or mismatched comment blocks (<!-- + define --> / <!-- + end -->,
// <!-- | define --> / <!-- | end -->, <!-- * define --> / <!-- * end -->).
//
// Usage: go run scripts/validate.go --src=path/to/html/dir
// Exit 0 if all files pass; 1 if any error (with file:line messages).
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reOpen  = regexp.MustCompile(`<!--\s*([*+|])\s+define\s+"([^"]*)"\s*-->`)
	reClose = regexp.MustCompile(`<!--\s*([*+|])\s+end\s*-->`)
)

func main() {
	src := flag.String("src", "", "directory containing .html files to validate")
	flag.Parse()
	if *src == "" {
		fmt.Fprintf(os.Stderr, "Usage: go run scripts/validate.go --src=DIR\n")
		os.Exit(2)
	}

	var hadErr bool
	err := filepath.Walk(*src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: read: %v\n", path, err)
			hadErr = true
			return nil
		}
		lines := strings.Split(string(content), "\n")
		var stack []string // e.g. "+", "|", "*"
		for i, line := range lines {
			lineNum := i + 1
			for _, m := range reClose.FindAllStringSubmatch(line, -1) {
				if len(m) < 2 {
					continue
				}
				delim := m[1]
				if len(stack) == 0 {
					fmt.Fprintf(os.Stderr, "%s:%d: unexpected <!-- %s end --> (no open block)\n", path, lineNum, delim)
					hadErr = true
					continue
				}
				top := stack[len(stack)-1]
				if top != delim {
					fmt.Fprintf(os.Stderr, "%s:%d: <!-- %s end --> does not match open <!-- %s define -->\n", path, lineNum, delim, top)
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
			fmt.Fprintf(os.Stderr, "%s:%d: unclosed block(s) <!-- %s define --> (missing <!-- %s end -->)\n", path, len(lines), strings.Join(stack, ", "), stack[len(stack)-1])
			hadErr = true
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "walk: %v\n", err)
		os.Exit(1)
	}
	if hadErr {
		os.Exit(1)
	}
}
