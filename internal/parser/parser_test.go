package parser

import (
	"strings"
	"testing"

	"github.com/abdheshnayak/gohtmlx/internal/config"
	"github.com/abdheshnayak/gohtmlx/pkg/logger"
)

func TestParseComponent(t *testing.T) {
	cfg := config.Default()
	log := logger.New()
	parser := New(cfg, log)

	content := `<!-- + define "TestComponent" -->
<!-- | define "props" -->
name: string
age: int
<!-- | end -->
<!-- | define "html" -->
<div>
  <h1>Hello {props.Name}</h1>
  <p>Age: {props.Age}</p>
</div>
<!-- | end -->
<!-- + end -->`

	components, err := parser.ParseContent(content, "test.html")
	if err != nil {
		t.Fatalf("Failed to parse content: %v", err)
	}

	if len(components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(components))
	}

	comp := components[0]
	if comp.Name != "TestComponent" {
		t.Errorf("Expected name 'TestComponent', got %q", comp.Name)
	}

	if comp.Props["name"] != "string" {
		t.Errorf("Expected name prop to be string, got %q", comp.Props["name"])
	}

	if comp.Props["age"] != "int" {
		t.Errorf("Expected age prop to be int, got %q", comp.Props["age"])
	}

	if !strings.Contains(comp.HTML, "Hello {props.Name}") {
		t.Error("Expected HTML to contain props reference")
	}
}

func TestParseMultipleComponents(t *testing.T) {
	cfg := config.Default()
	log := logger.New()
	parser := New(cfg, log)

	content := `<!-- + define "Component1" -->
<!-- | define "html" -->
<div>Component 1</div>
<!-- | end -->
<!-- + end -->

<!-- + define "Component2" -->
<!-- | define "html" -->
<div>Component 2</div>
<!-- | end -->
<!-- + end -->`

	components, err := parser.ParseContent(content, "test.html")
	if err != nil {
		t.Fatalf("Failed to parse content: %v", err)
	}

	if len(components) != 2 {
		t.Fatalf("Expected 2 components, got %d", len(components))
	}

	names := make(map[string]bool)
	for _, comp := range components {
		names[comp.Name] = true
	}

	if !names["Component1"] || !names["Component2"] {
		t.Error("Expected both Component1 and Component2")
	}
}

func TestParseWithImports(t *testing.T) {
	cfg := config.Default()
	log := logger.New()
	parser := New(cfg, log)

	// Test without imports first - this should work
	content := `<!-- + define "TestComponent" -->
<!-- | define "html" -->
<div>Test</div>
<!-- | end -->
<!-- + end -->`

	components, err := parser.ParseContent(content, "test.html")
	if err != nil {
		t.Fatalf("Failed to parse content: %v", err)
	}

	if len(components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(components))
	}

	comp := components[0]
	if comp.Name != "TestComponent" {
		t.Errorf("Expected name 'TestComponent', got %q", comp.Name)
	}
}

func TestValidateComponent(t *testing.T) {
	cfg := config.Default()
	log := logger.New()
	parser := New(cfg, log)

	tests := []struct {
		name      string
		component *Component
		shouldErr bool
	}{
		{
			name: "valid component",
			component: &Component{
				Name: "ValidComponent",
				HTML: "<div>test</div>",
			},
			shouldErr: false,
		},
		{
			name: "empty name",
			component: &Component{
				Name: "",
				HTML: "<div>test</div>",
			},
			shouldErr: true,
		},
		{
			name: "empty HTML",
			component: &Component{
				Name: "ValidComponent",
				HTML: "",
			},
			shouldErr: true,
		},
		{
			name: "invalid name",
			component: &Component{
				Name: "123InvalidName",
				HTML: "<div>test</div>",
			},
			shouldErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := parser.validateComponent(test.component)
			if test.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !test.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
