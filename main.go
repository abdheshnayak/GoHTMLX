package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/abdheshnayak/gohtmlx/pkg/transpiler"
	"github.com/abdheshnayak/gohtmlx/pkg/utils"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of gohtmlx:\n")
		fmt.Fprintf(os.Stderr, "  gohtmlx --src=DIR --dist=DIR\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExit codes:\n")
		fmt.Fprintf(os.Stderr, "  0  success\n")
		fmt.Fprintf(os.Stderr, "  1  transpilation failed (parse, codegen, or write error)\n")
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
