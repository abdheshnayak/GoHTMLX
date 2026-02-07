package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/abdheshnayak/gohtmlx/pkg/transpiler"
	"github.com/abdheshnayak/gohtmlx/pkg/utils"
	"github.com/abdheshnayak/gohtmlx/pkg/validate"
)

func main() {
	if len(os.Args) >= 2 && (os.Args[1] == "validate" || os.Args[1] == "check") {
		runValidate(os.Args[2:])
		return
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of gohtmlx:\n")
		fmt.Fprintf(os.Stderr, "  gohtmlx --src=DIR --dist=DIR     transpile .html components to Go\n")
		fmt.Fprintf(os.Stderr, "  gohtmlx validate --src=DIR       check comment structure (unclosed define/end)\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExit codes:\n")
		fmt.Fprintf(os.Stderr, "  0  success\n")
		fmt.Fprintf(os.Stderr, "  1  transpilation/validation failed\n")
		fmt.Fprintf(os.Stderr, "  2  invalid arguments or missing flags\n")
	}

	src := flag.String("src", "", "source directory containing .html components")
	dist := flag.String("dist", "", "destination directory for generated Go code")
	singleFile := flag.Bool("single-file", false, "emit one comp_generated.go (legacy); default is one file per component")
	pkg := flag.String("pkg", "gohtmlxc", "generated package name")
	flag.Parse()

	if *src == "" || *dist == "" {
		flag.Usage()
		os.Exit(2)
	}

	utils.Log = utils.NewSlogLogger(slog.Default())
	opts := &transpiler.RunOptions{SingleFile: *singleFile, Pkg: *pkg}
	if err := transpiler.Run(*src, *dist, opts); err != nil {
		utils.Log.Error("transpiling failed", "err", err)
		os.Exit(1)
	}
}

func runValidate(args []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	src := fs.String("src", "", "directory containing .html files to validate")
	_ = fs.Parse(args)
	if *src == "" {
		fmt.Fprintf(os.Stderr, "Usage: gohtmlx validate --src=DIR\n")
		os.Exit(2)
	}
	if err := validate.Run(*src, os.Stderr); err != nil {
		os.Exit(1)
	}
}
