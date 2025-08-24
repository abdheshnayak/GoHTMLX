package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration for GoHTMLX
type Config struct {
	SourceDir   string            `yaml:"source_dir"`
	OutputDir   string            `yaml:"output_dir"`
	PackageName string            `yaml:"package_name"`
	LogLevel    string            `yaml:"log_level"`
	Watch       WatchConfig       `yaml:"watch"`
	Template    TemplateConfig    `yaml:"template"`
	Generation  GenerationConfig  `yaml:"generation"`
	Extensions  map[string]string `yaml:"extensions"`
}

// WatchConfig holds watch-specific configuration
type WatchConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Extensions  []string `yaml:"extensions"`
	IgnoreFiles []string `yaml:"ignore_files"`
	Debounce    int      `yaml:"debounce_ms"`
}

// TemplateConfig holds template parsing configuration
type TemplateConfig struct {
	StartDelim string `yaml:"start_delim"`
	EndDelim   string `yaml:"end_delim"`
	PropDelim  string `yaml:"prop_delim"`
	PropEnd    string `yaml:"prop_end"`
}

// GenerationConfig holds code generation configuration
type GenerationConfig struct {
	FormatCode     bool   `yaml:"format_code"`
	AddComments    bool   `yaml:"add_comments"`
	ImportPath     string `yaml:"import_path"`
	ElementPackage string `yaml:"element_package"`
}

// Default returns a default configuration
func Default() *Config {
	return &Config{
		SourceDir:   "src",
		OutputDir:   "dist",
		PackageName: "gohtmlxc",
		LogLevel:    "info",
		Watch: WatchConfig{
			Enabled:     true,
			Extensions:  []string{".html", ".htm"},
			IgnoreFiles: []string{"*.tmp", ".*"},
			Debounce:    300,
		},
		Template: TemplateConfig{
			StartDelim: "<!-- +",
			EndDelim:   " -->",
			PropDelim:  "<!-- |",
			PropEnd:    " -->",
		},
		Generation: GenerationConfig{
			FormatCode:     true,
			AddComments:    true,
			ImportPath:     "github.com/abdheshnayak/gohtmlx/pkg/element",
			ElementPackage: "element",
		},
		Extensions: map[string]string{
			".html": "html",
			".htm":  "html",
		},
	}
}

// LoadFromFile loads configuration from a YAML file
func (c *Config) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// SaveToFile saves configuration to a YAML file
func (c *Config) SaveToFile(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.SourceDir == "" {
		return fmt.Errorf("source directory is required")
	}

	if c.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}

	if c.PackageName == "" {
		return fmt.Errorf("package name is required")
	}

	// Check if source directory exists
	if _, err := os.Stat(c.SourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", c.SourceDir)
	}

	// Validate log level
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	return nil
}
