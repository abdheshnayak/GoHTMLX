package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/abdheshnayak/gohtmlx/internal/compiler"
	"github.com/abdheshnayak/gohtmlx/internal/config"
	"github.com/abdheshnayak/gohtmlx/internal/watcher"
	"github.com/abdheshnayak/gohtmlx/pkg/logger"
)

// App represents the CLI application
type App struct {
	Version string
	Commit  string
	Date    string
	logger  logger.Logger
}

// NewApp creates a new CLI application instance
func NewApp() *App {
	return &App{
		logger: logger.New(),
	}
}

// Run executes the CLI application
func (a *App) Run(args []string) error {
	if len(args) < 2 {
		return a.showHelp()
	}

	command := args[1]

	switch command {
	case "build":
		return a.runBuild(args[2:])
	case "watch":
		return a.runWatch(args[2:])
	case "version", "-v", "--version":
		return a.showVersion()
	case "help", "-h", "--help":
		return a.showHelp()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (a *App) runBuild(args []string) error {
	cfg, err := a.parseFlags(args)
	if err != nil {
		return err
	}

	compiler := compiler.New(cfg, a.logger)
	return compiler.Build()
}

func (a *App) runWatch(args []string) error {
	cfg, err := a.parseFlags(args)
	if err != nil {
		return err
	}

	w := watcher.New(cfg, a.logger)
	return w.Start()
}

func (a *App) parseFlags(args []string) (*config.Config, error) {
	fs := flag.NewFlagSet("gohtmlx", flag.ExitOnError)

	var (
		srcDir     = fs.String("src", "", "Source directory containing HTML components")
		distDir    = fs.String("dist", "", "Output directory for generated Go code")
		configFile = fs.String("config", "", "Configuration file path")
		verbose    = fs.Bool("verbose", false, "Enable verbose logging")
		quiet      = fs.Bool("quiet", false, "Suppress output")
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// Load configuration
	cfg := config.Default()
	if *configFile != "" {
		if err := cfg.LoadFromFile(*configFile); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// Override with command line flags
	if *srcDir != "" {
		cfg.SourceDir = *srcDir
	}
	if *distDir != "" {
		cfg.OutputDir = *distDir
	}
	if *verbose {
		cfg.LogLevel = "debug"
	}
	if *quiet {
		cfg.LogLevel = "error"
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Convert to absolute paths
	if absPath, err := filepath.Abs(cfg.SourceDir); err == nil {
		cfg.SourceDir = absPath
	}
	if absPath, err := filepath.Abs(cfg.OutputDir); err == nil {
		cfg.OutputDir = absPath
	}

	return cfg, nil
}

func (a *App) showVersion() error {
	fmt.Printf("gohtmlx version %s\n", a.Version)
	fmt.Printf("Commit: %s\n", a.Commit)
	fmt.Printf("Built: %s\n", a.Date)
	return nil
}

func (a *App) showHelp() error {
	fmt.Fprintf(os.Stderr, `GoHTMLX - HTML Components with Go

USAGE:
    gohtmlx <command> [options]

COMMANDS:
    build       Compile HTML components to Go code
    watch       Watch for changes and rebuild automatically
    version     Show version information
    help        Show this help message

BUILD OPTIONS:
    --src <dir>      Source directory containing HTML components
    --dist <dir>     Output directory for generated Go code
    --config <file>  Configuration file path
    --verbose        Enable verbose logging
    --quiet          Suppress output

EXAMPLES:
    gohtmlx build --src ./components --dist ./generated
    gohtmlx watch --src ./components --dist ./generated --verbose
    gohtmlx build --config gohtmlx.config.yaml

For more information, visit: https://github.com/abdheshnayak/gohtmlx
`)
	return nil
}
