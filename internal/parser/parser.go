package parser

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/abdheshnayak/gohtmlx/internal/config"
	"github.com/abdheshnayak/gohtmlx/pkg/logger"
	"sigs.k8s.io/yaml"
)

// Component represents a parsed HTML component
type Component struct {
	Name       string            `json:"name"`
	Props      map[string]string `json:"props"`
	HTML       string            `json:"html"`
	FilePath   string            `json:"file_path"`
	Imports    []string          `json:"imports"`
	Attributes map[string]string `json:"attributes"`
}

// Parser handles parsing of HTML component files
type Parser struct {
	config *config.Config
	logger logger.Logger
}

// New creates a new parser instance
func New(cfg *config.Config, log logger.Logger) *Parser {
	return &Parser{
		config: cfg,
		logger: log,
	}
}

// ParseFiles parses multiple HTML component files
func (p *Parser) ParseFiles(filePaths []string) ([]*Component, error) {
	var components []*Component

	for _, filePath := range filePaths {
		fileComponents, err := p.ParseFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
		}
		components = append(components, fileComponents...)
	}

	return components, nil
}

// ParseFile parses a single HTML component file
func (p *Parser) ParseFile(filePath string) ([]*Component, error) {
	p.logger.Debug("Parsing file", "file", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return p.ParseContent(string(content), filePath)
}

// ParseContent parses HTML component content
func (p *Parser) ParseContent(content, filePath string) ([]*Component, error) {
	// First, extract global imports
	imports, err := p.extractImports(content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract imports: %w", err)
	}

	// Parse component sections
	components, err := p.parseComponents(content, filePath, imports)
	if err != nil {
		return nil, fmt.Errorf("failed to parse components: %w", err)
	}

	return components, nil
}

// extractImports extracts global imports from the content
func (p *Parser) extractImports(content string) ([]string, error) {
	tmpl, err := template.New("parser_imports").Delims("<!-- *", " -->").Parse(content)
	if err != nil {
		return nil, err
	}

	sections, err := p.parseSections(tmpl)
	if err != nil {
		return nil, err
	}

	var imports []string
	if importsStr, ok := sections["imports"]; ok {
		for _, imp := range strings.Split(strings.TrimSpace(importsStr), "\n") {
			imp = strings.TrimSpace(imp)
			if imp != "" {
				imports = append(imports, imp)
			}
		}
	}

	return imports, nil
}

// parseComponents parses individual components from content
func (p *Parser) parseComponents(content, filePath string, globalImports []string) ([]*Component, error) {
	tmpl, err := template.New("components").Delims(p.config.Template.StartDelim, p.config.Template.EndDelim).Parse(content)
	if err != nil {
		return nil, err
	}

	sections, err := p.parseSections(tmpl)
	if err != nil {
		return nil, err
	}

	var components []*Component
	for name, sectionContent := range sections {
		if name == "imports" {
			continue // Skip imports section
		}

		component, err := p.parseComponent(name, sectionContent, filePath, globalImports)
		if err != nil {
			return nil, fmt.Errorf("failed to parse component %s: %w", name, err)
		}

		components = append(components, component)
	}

	return components, nil
}

// parseComponent parses a single component section
func (p *Parser) parseComponent(name, content, filePath string, globalImports []string) (*Component, error) {
	component := &Component{
		Name:     name,
		FilePath: filePath,
		Props:    make(map[string]string),
		Imports:  globalImports,
	}

	// Parse component subsections (props, html, etc.)
	tmpl, err := template.New("component").Delims(p.config.Template.PropDelim, p.config.Template.PropEnd).Parse(content)
	if err != nil {
		return nil, err
	}

	subsections, err := p.parseSections(tmpl)
	if err != nil {
		return nil, err
	}

	// Extract props
	if propsStr, ok := subsections["props"]; ok {
		if err := yaml.Unmarshal([]byte(propsStr), &component.Props); err != nil {
			return nil, fmt.Errorf("failed to parse props: %w", err)
		}
	}

	// Extract HTML
	if htmlStr, ok := subsections["html"]; ok {
		component.HTML = strings.TrimSpace(htmlStr)
	}

	// Validate component
	if err := p.validateComponent(component); err != nil {
		return nil, fmt.Errorf("component validation failed: %w", err)
	}

	return component, nil
}

// parseSections parses template sections
func (p *Parser) parseSections(tmpl *template.Template) (map[string]string, error) {
	sections := make(map[string]string)

	for _, t := range tmpl.Templates() {
		if t.Name() == tmpl.Name() {
			continue // Skip the main template
		}

		var buf strings.Builder
		if err := t.Execute(&buf, nil); err != nil {
			return nil, err
		}
		sections[t.Name()] = buf.String()
	}

	return sections, nil
}

// validateComponent validates a parsed component
func (p *Parser) validateComponent(component *Component) error {
	if component.Name == "" {
		return fmt.Errorf("component name is required")
	}

	if component.HTML == "" {
		return fmt.Errorf("component HTML is required")
	}

	// Validate component name (should be valid Go identifier)
	if !isValidIdentifier(component.Name) {
		return fmt.Errorf("invalid component name: %s", component.Name)
	}

	return nil
}

// isValidIdentifier checks if a string is a valid Go identifier
func isValidIdentifier(name string) bool {
	if len(name) == 0 {
		return false
	}

	// Must start with letter or underscore
	first := name[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_') {
		return false
	}

	// Rest must be letters, digits, or underscores
	for i := 1; i < len(name); i++ {
		char := name[i]
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}

	return true
}
