// validate is a small checker for GoHTMLX .html component files. It reports
// unclosed or mismatched comment blocks (<!-- + define --> / <!-- + end -->,
// <!-- | define --> / <!-- | end -->, <!-- * define --> / <!-- * end -->).
//
// Usage: go run scripts/validate.go --src=path/to/html/dir
// Or:   gohtmlx validate --src=path/to/html/dir
// Exit 0 if all files pass; 1 if any error (with file:line messages).
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/abdheshnayak/gohtmlx/pkg/validate"
)

func main() {
	src := flag.String("src", "", "directory containing .html files to validate")
	flag.Parse()
	if *src == "" {
		fmt.Fprintf(os.Stderr, "Usage: go run scripts/validate.go --src=DIR\n")
		os.Exit(2)
	}

	if err := validate.Run(*src, os.Stderr); err != nil {
		os.Exit(1)
	}
}
