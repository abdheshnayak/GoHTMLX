package element

import (
	"strings"
	"testing"
)

func TestNewHtml_SimpleElement(t *testing.T) {
	h, err := NewHtml([]byte("<div>hello</div>"))
	if err != nil {
		t.Fatalf("NewHtml: %v", err)
	}
	comps := map[string]CompInfo{}
	out, err := h.RenderGolangCode(comps)
	if err != nil {
		t.Fatalf("RenderGolangCode: %v", err)
	}
	if !strings.Contains(out, "E(") || !strings.Contains(out, "div") {
		t.Errorf("expected E( and div in output, got: %s", out)
	}
}

func TestNewHtml_WithPropsExpression(t *testing.T) {
	h, err := NewHtml([]byte("<span>{props.Name}</span>"))
	if err != nil {
		t.Fatalf("NewHtml: %v", err)
	}
	comps := map[string]CompInfo{}
	out, err := h.RenderGolangCode(comps)
	if err != nil {
		t.Fatalf("RenderGolangCode: %v", err)
	}
	if !strings.Contains(out, "Name") {
		t.Errorf("expected props.Name in output, got: %s", out)
	}
}

func TestNewHtml_CustomComponent(t *testing.T) {
	// Custom tag must be in comps; HTML parser lowercases so key is "greet"
	comps := map[string]CompInfo{
		"greet": {Name: "Greet", Props: map[string]string{"name": "Name"}},
	}
	h, err := NewHtml([]byte("<greet name={props.UserName}></greet>"))
	if err != nil {
		t.Fatalf("NewHtml: %v", err)
	}
	out, err := h.RenderGolangCode(comps)
	if err != nil {
		t.Fatalf("RenderGolangCode: %v", err)
	}
	if !strings.Contains(out, "GreetComp(") {
		t.Errorf("expected GreetComp in output, got: %s", out)
	}
}

func TestNewHtml_ForElement(t *testing.T) {
	h, err := NewHtml([]byte(`<for items={props.Items} as="item"><span>{item}</span></for>`))
	if err != nil {
		t.Fatalf("NewHtml: %v", err)
	}
	comps := map[string]CompInfo{}
	out, err := h.RenderGolangCode(comps)
	if err != nil {
		t.Fatalf("RenderGolangCode: %v", err)
	}
	if !strings.Contains(out, "for _,") || !strings.Contains(out, "range") {
		t.Errorf("expected for/range in output, got: %s", out)
	}
}

func TestNewHtml_InvalidHTML(t *testing.T) {
	// html.ParseFragment can be lenient; test that we don't panic and get some output or error
	h, err := NewHtml([]byte("<div>ok</div>"))
	if err != nil {
		t.Fatalf("valid HTML should parse: %v", err)
	}
	_ = h
	// Severely broken HTML may or may not error depending on parser
	_, _ = NewHtml([]byte("<<<"))
}
