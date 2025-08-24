package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abdheshnayak/gohtmlx/internal/config"
	"github.com/abdheshnayak/gohtmlx/internal/generator"
	"github.com/abdheshnayak/gohtmlx/internal/parser"
	"github.com/abdheshnayak/gohtmlx/pkg/logger"
)

// Compiler handles the compilation process
type Compiler struct {
	config    *config.Config
	logger    logger.Logger
	parser    *parser.Parser
	generator *generator.Generator
}

// New creates a new compiler instance
func New(cfg *config.Config, log logger.Logger) *Compiler {
	return &Compiler{
		config:    cfg,
		logger:    log,
		parser:    parser.New(cfg, log),
		generator: generator.New(cfg, log),
	}
}

// Build compiles all HTML components to Go code
func (c *Compiler) Build() error {
	start := time.Now()
	c.logger.Info("Starting compilation...", "src", c.config.SourceDir, "dist", c.config.OutputDir)

	// Find all HTML files
	files, err := c.findHTMLFiles()
	if err != nil {
		return fmt.Errorf("failed to find HTML files: %w", err)
	}

	if len(files) == 0 {
		c.logger.Warn("No HTML files found", "src", c.config.SourceDir)
		return nil
	}

	c.logger.Debug("Found HTML files", "count", len(files), "files", files)

	// Parse all files
	components, err := c.parser.ParseFiles(files)
	if err != nil {
		return fmt.Errorf("failed to parse files: %w", err)
	}

	c.logger.Debug("Parsed components", "count", len(components))

	// Generate Go code
	if err := c.generator.Generate(components); err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	duration := time.Since(start)
	c.logger.Info("Compilation completed", "duration", duration, "components", len(components))

	return nil
}

// findHTMLFiles recursively finds all HTML files in the source directory
func (c *Compiler) findHTMLFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(c.config.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file has a valid extension
		ext := strings.ToLower(filepath.Ext(path))
		if _, ok := c.config.Extensions[ext]; ok {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// BuildFile compiles a single HTML file
func (c *Compiler) BuildFile(filePath string) error {
	c.logger.Debug("Building single file", "file", filePath)

	// Parse the file
	components, err := c.parser.ParseFiles([]string{filePath})
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// Generate Go code
	if err := c.generator.Generate(components); err != nil {
		return fmt.Errorf("failed to generate code for file %s: %w", filePath, err)
	}

	c.logger.Debug("File built successfully", "file", filePath, "components", len(components))
	return nil
}
