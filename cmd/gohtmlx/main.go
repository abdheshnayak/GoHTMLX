package main

import (
	"fmt"
	"os"

	"github.com/abdheshnayak/gohtmlx/internal/cli"
)

var (
	// Version information (set during build)
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Commit = Commit
	app.Date = Date

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
