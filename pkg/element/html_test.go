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

func TestNewHtml_IfElement(t *testing.T) {
	h, err := NewHtml([]byte(`<if condition={props.Show}><span>yes</span></if>`))
	if err != nil {
		t.Fatalf("NewHtml: %v", err)
	}
	comps := map[string]CompInfo{}
	out, err := h.RenderGolangCode(comps)
	if err != nil {
		t.Fatalf("RenderGolangCode: %v", err)
	}
	if !strings.Contains(out, "if ") || !strings.Contains(out, "return []Element{") {
		t.Errorf("expected if and return []Element in output, got: %s", out)
	}
	if !strings.Contains(out, "Show") {
		t.Errorf("expected condition (Show) in output, got: %s", out)
	}
}

func TestNewHtml_IfElseElement(t *testing.T) {
	h, err := NewHtml([]byte(`<if condition={props.A}><span>a</span></if><else><span>no</span></else>`))
	if err != nil {
		t.Fatalf("NewHtml: %v", err)
	}
	comps := map[string]CompInfo{}
	out, err := h.RenderGolangCode(comps)
	if err != nil {
		t.Fatalf("RenderGolangCode: %v", err)
	}
	if !strings.Contains(out, "if ") || !strings.Contains(out, "return []Element{") {
		t.Errorf("expected if and return in output, got: %s", out)
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

func TestSlotNamesFromHTML(t *testing.T) {
	names, err := SlotNamesFromHTML([]byte(`<div><slot name="header"></slot><slot name="footer"/></div>`))
	if err != nil {
		t.Fatalf("SlotNamesFromHTML: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 slot names, got %d: %v", len(names), names)
	}
	// Order may vary; check both are present
	seen := make(map[string]bool)
	for _, n := range names {
		seen[n] = true
	}
	if !seen["header"] || !seen["footer"] {
		t.Errorf("expected header and footer, got %v", names)
	}
}

func TestNewHtml_SlotPlaceholder(t *testing.T) {
	// Layout template: uses <slot name="header"/> placeholder
	comps := map[string]CompInfo{
		"layout": {Name: "Layout", Props: map[string]string{"slotheader": "slotHeader"}},
	}
	h, err := NewHtml([]byte(`<div><slot name="header"/></div>`))
	if err != nil {
		t.Fatalf("NewHtml: %v", err)
	}
	out, err := h.RenderGolangCode(comps)
	if err != nil {
		t.Fatalf("RenderGolangCode: %v", err)
	}
	if !strings.Contains(out, "props.SlotHeader") {
		t.Errorf("expected props.SlotHeader in output, got: %s", out)
	}
}

func TestNewHtml_SlotContentFromCaller(t *testing.T) {
	// Caller passes <slot name="header">content</slot> into Layout
	comps := map[string]CompInfo{
		"layout": {Name: "Layout", Props: map[string]string{"slotheader": "slotHeader"}},
	}
	h, err := NewHtml([]byte(`<layout><slot name="header"><span>title</span></slot></layout>`))
	if err != nil {
		t.Fatalf("NewHtml: %v", err)
	}
	out, err := h.RenderGolangCode(comps)
	if err != nil {
		t.Fatalf("RenderGolangCode: %v", err)
	}
	if !strings.Contains(out, "LayoutComp(") || !strings.Contains(out, "SlotHeader") {
		t.Errorf("expected LayoutComp and SlotHeader in output, got: %s", out)
	}
	if !strings.Contains(out, "span") {
		t.Errorf("expected slot content (span) in output, got: %s", out)
	}
}

func TestNewHtml_MultipleBracedExpressions(t *testing.T) {
	// Text with multiple {expr} must produce valid Go (comma-separated args), not one broken expression
	h, err := NewHtml([]byte("<p>{props.Author} — {props.Role}</p>"))
	if err != nil {
		t.Fatalf("NewHtml: %v", err)
	}
	comps := map[string]CompInfo{}
	out, err := h.RenderGolangCode(comps)
	if err != nil {
		t.Fatalf("RenderGolangCode: %v", err)
	}
	// Should have both identifiers and a comma (R(..., ..., ...))
	if !strings.Contains(out, "props.Author") || !strings.Contains(out, "props.Role") {
		t.Errorf("expected both props.Author and props.Role in output, got: %s", out)
	}
	if !strings.Contains(out, ",") {
		t.Errorf("multiple expressions should be comma-separated in R(...), got: %s", out)
	}
}

// FuzzProcessRaws exercises processRaws with arbitrary inputs to catch panics or invalid output.
// Run with: go test -fuzz=FuzzProcessRaws -fuzztime=30s ./pkg/element/
func FuzzProcessRaws(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1e5 {
			t.Skip("input too large")
		}
		_ = processRaws(string(data))
	})
}
