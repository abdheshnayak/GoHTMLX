package generator

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/abdheshnayak/gohtmlx/internal/config"
	"github.com/abdheshnayak/gohtmlx/internal/parser"
	"github.com/abdheshnayak/gohtmlx/pkg/element"
	"github.com/abdheshnayak/gohtmlx/pkg/logger"
)

// Generator handles Go code generation from parsed components
type Generator struct {
	config *config.Config
	logger logger.Logger
}

// New creates a new generator instance
func New(cfg *config.Config, log logger.Logger) *Generator {
	return &Generator{
		config: cfg,
		logger: log,
	}
}

// Generate generates Go code from parsed components
func (g *Generator) Generate(components []*parser.Component) error {
	if len(components) == 0 {
		g.logger.Warn("No components to generate")
		return nil
	}

	g.logger.Debug("Generating Go code", "components", len(components))

	// Create component info map for cross-referencing
	compInfo := make(map[string]element.CompInfo)
	for _, comp := range components {
		compInfo[strings.ToLower(comp.Name)] = element.CompInfo{
			Name:  comp.Name,
			Props: g.normalizeProps(comp.Props),
		}
	}

	// Generate structs and functions
	structs, functions, err := g.generateCode(components, compInfo)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// Collect all imports
	imports := g.collectImports(components)

	// Generate final source code
	source, err := g.generateSourceFile(structs, functions, imports)
	if err != nil {
		return fmt.Errorf("failed to generate source file: %w", err)
	}

	// Write to output file
	outputPath := filepath.Join(g.config.OutputDir, g.config.PackageName, "components_generated.go")
	if err := g.writeSourceFile(outputPath, source); err != nil {
		return fmt.Errorf("failed to write source file: %w", err)
	}

	g.logger.Info("Go code generated successfully", "output", outputPath)
	return nil
}

// generateCode generates structs and functions for components
func (g *Generator) generateCode(components []*parser.Component, compInfo map[string]element.CompInfo) ([]string, []string, error) {
	var structs []string
	var functions []string

	for _, comp := range components {
		// Generate struct
		structCode, err := g.generateStruct(comp)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate struct for %s: %w", comp.Name, err)
		}
		structs = append(structs, structCode)

		// Generate HTML rendering function
		funcCode, err := g.generateFunction(comp, compInfo)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate function for %s: %w", comp.Name, err)
		}
		functions = append(functions, funcCode)
	}

	return structs, functions, nil
}

// generateStruct generates a Go struct for a component
func (g *Generator) generateStruct(comp *parser.Component) (string, error) {
	var buf strings.Builder

	if g.config.Generation.AddComments {
		buf.WriteString(fmt.Sprintf("// %s represents the props for the %s component\n", comp.Name, comp.Name))
	}

	buf.WriteString(fmt.Sprintf("type %s struct {\n", comp.Name))

	// Add props as struct fields
	for propName, propType := range comp.Props {
		fieldName := g.capitalize(propName)
		buf.WriteString(fmt.Sprintf("\t%s %s\n", fieldName, propType))
	}

	// Add Attrs field
	buf.WriteString("\tAttrs element.Attrs\n")
	buf.WriteString("}\n")

	if g.config.Generation.FormatCode {
		formatted, err := format.Source([]byte(buf.String()))
		if err != nil {
			g.logger.Warn("Failed to format struct", "component", comp.Name, "error", err)
			return buf.String(), nil
		}
		return string(formatted), nil
	}

	return buf.String(), nil
}

// generateFunction generates a Go function for a component
func (g *Generator) generateFunction(comp *parser.Component, compInfo map[string]element.CompInfo) (string, error) {
	// Parse HTML and generate element code
	html, err := element.NewHtml([]byte(comp.HTML))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	elementCode, err := html.RenderGolangCode(compInfo)
	if err != nil {
		return "", fmt.Errorf("failed to render Go code: %w", err)
	}

	var buf strings.Builder

	if g.config.Generation.AddComments {
		buf.WriteString(fmt.Sprintf("// %sComp creates a %s component with the given props and attributes\n", comp.Name, comp.Name))
	}

	// Component function
	buf.WriteString(fmt.Sprintf("func %sComp(props %s, attrs element.Attrs, children ...element.Element) element.Element {\n", comp.Name, comp.Name))
	buf.WriteString("\tprops.Attrs = attrs\n")
	buf.WriteString("\tif props.Attrs == nil {\n")
	buf.WriteString("\t\tprops.Attrs = element.Attrs{}\n")
	buf.WriteString("\t}\n")
	buf.WriteString(fmt.Sprintf("\treturn %s\n", elementCode))
	buf.WriteString("}\n\n")

	// Get method for convenient usage
	if g.config.Generation.AddComments {
		buf.WriteString(fmt.Sprintf("// Get returns the %s component as an Element\n", comp.Name))
	}
	buf.WriteString(fmt.Sprintf("func (c %s) Get(children ...element.Element) element.Element {\n", comp.Name))
	buf.WriteString(fmt.Sprintf("\treturn %sComp(c, c.Attrs, children...)\n", comp.Name))
	buf.WriteString("}\n")

	return buf.String(), nil
}

// generateSourceFile generates the complete Go source file
func (g *Generator) generateSourceFile(structs, functions, imports []string) (string, error) {
	tmplStr := `package {{.PackageName}}

import (
	"{{.ElementPackage}}"
{{range .Imports}}	{{.}}
{{end}})

{{range .Structs}}{{.}}

{{end}}{{range .Functions}}{{.}}

{{end}}`

	tmpl, err := template.New("source").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	data := struct {
		PackageName    string
		ElementPackage string
		Imports        []string
		Structs        []string
		Functions      []string
	}{
		PackageName:    g.config.PackageName,
		ElementPackage: g.config.Generation.ImportPath,
		Imports:        imports,
		Structs:        structs,
		Functions:      functions,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	source := buf.String()

	if g.config.Generation.FormatCode {
		formatted, err := format.Source([]byte(source))
		if err != nil {
			g.logger.Warn("Failed to format generated source", "error", err)
			return source, nil
		}
		return string(formatted), nil
	}

	return source, nil
}

// writeSourceFile writes the generated source code to a file
func (g *Generator) writeSourceFile(outputPath, source string) error {
	// Create output directory
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write source file
	if err := os.WriteFile(outputPath, []byte(source), 0644); err != nil {
		return fmt.Errorf("failed to write source file: %w", err)
	}

	return nil
}

// collectImports collects all unique imports from components
func (g *Generator) collectImports(components []*parser.Component) []string {
	importSet := make(map[string]bool)

	for _, comp := range components {
		for _, imp := range comp.Imports {
			importSet[strings.TrimSpace(imp)] = true
		}
	}

	var imports []string
	for imp := range importSet {
		if imp != "" {
			imports = append(imports, imp)
		}
	}

	return imports
}

// normalizeProps normalizes prop names for component info
// Maps HTML attribute names to Go field names
func (g *Generator) normalizeProps(props map[string]string) map[string]string {
	normalized := make(map[string]string)
	for propName, _ := range props {
		// Map lowercase HTML attribute to capitalized Go field name
		htmlAttr := strings.ToLower(propName)
		goFieldName := g.capitalize(propName)
		normalized[htmlAttr] = goFieldName
	}
	return normalized
}

// capitalize capitalizes the first letter of a string
func (g *Generator) capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
