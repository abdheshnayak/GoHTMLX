package element

import (
	"strings"
	"testing"
)

func TestBasicElement(t *testing.T) {
	elem := Div(Attrs{"class": "container"}, Text("Hello World"))

	var buf strings.Builder
	_, err := elem.Render(&buf)
	if err != nil {
		t.Fatalf("Failed to render element: %v", err)
	}

	expected := `<div class="container">Hello World</div>`
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestSelfClosingTags(t *testing.T) {
	elem := Img("test.jpg", "Test Image", Attrs{"width": "100"})

	var buf strings.Builder
	_, err := elem.Render(&buf)
	if err != nil {
		t.Fatalf("Failed to render element: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, `src="test.jpg"`) {
		t.Error("Expected src attribute")
	}
	if !strings.Contains(result, `alt="Test Image"`) {
		t.Error("Expected alt attribute")
	}
	if !strings.Contains(result, ` />`) {
		t.Error("Expected self-closing tag")
	}
}

func TestClassNames(t *testing.T) {
	tests := []struct {
		input    []any
		expected string
	}{
		{[]any{"foo", "bar"}, "foo bar"},
		{[]any{"foo", "", "bar"}, "foo bar"},
		{[]any{[]string{"foo", "bar"}}, "foo bar"},
		{[]any{map[string]bool{"foo": true, "bar": false}}, "foo"},
	}

	for _, test := range tests {
		result := ClassNames(test.input...)
		if result != test.expected {
			t.Errorf("ClassNames(%v) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestAttributeHelpers(t *testing.T) {
	attrs := MergeAttrs(
		ID("test-id"),
		Class("foo", "bar"),
		DataAttr("value", "123"),
	)

	if attrs["id"] != "test-id" {
		t.Error("Expected id attribute")
	}
	if attrs["class"] != "foo bar" {
		t.Error("Expected class attribute")
	}
	if attrs["data-value"] != "123" {
		t.Error("Expected data-value attribute")
	}
}

func TestConditionalRendering(t *testing.T) {
	elem1 := If(true, Text("Shown"))
	elem2 := If(false, Text("Hidden"))

	var buf1, buf2 strings.Builder
	elem1.Render(&buf1)
	elem2.Render(&buf2)

	if buf1.String() != "Shown" {
		t.Error("Expected conditional element to render when true")
	}
	if buf2.String() != "" {
		t.Error("Expected conditional element to not render when false")
	}
}

func TestForLoop(t *testing.T) {
	items := []string{"a", "b", "c"}
	elem := For(items, func(item string, index int) Element {
		return Li(nil, Text(item))
	})

	var buf strings.Builder
	elem.Render(&buf)

	result := buf.String()
	expected := "<li>a</li><li>b</li><li>c</li>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFragment(t *testing.T) {
	elem := Fragment(
		Text("Hello "),
		Strong(nil, Text("World")),
		Text("!"),
	)

	var buf strings.Builder
	elem.Render(&buf)

	expected := "Hello <strong>World</strong>!"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestHTMLEscaping(t *testing.T) {
	elem := Div(Attrs{"title": "<script>"}, Text("<script>alert('xss')</script>"))

	var buf strings.Builder
	elem.Render(&buf)

	result := buf.String()
	if strings.Contains(result, "<script>alert") {
		t.Error("HTML should be escaped in text content")
	}
	if !strings.Contains(result, `title="&lt;script&gt;"`) {
		t.Error("HTML should be escaped in attributes")
	}
}

func TestBooleanAttributes(t *testing.T) {
	elem := Input("checkbox", Attrs{"checked": true, "disabled": false})

	var buf strings.Builder
	elem.Render(&buf)

	result := buf.String()
	if !strings.Contains(result, "checked") {
		t.Error("Expected boolean attribute to be present when true")
	}
	if strings.Contains(result, "disabled") {
		t.Error("Expected boolean attribute to be absent when false")
	}
}

func TestNestedElements(t *testing.T) {
	elem := Div(Class("container"),
		Header(nil,
			H1(nil, Text("Title")),
		),
		Main(nil,
			P(nil, Text("Content")),
		),
	)

	var buf strings.Builder
	elem.Render(&buf)

	result := buf.String()
	if !strings.Contains(result, "<div class=\"container\">") {
		t.Error("Expected outer div")
	}
	if !strings.Contains(result, "<header>") {
		t.Error("Expected header element")
	}
	if !strings.Contains(result, "<h1>Title</h1>") {
		t.Error("Expected h1 element")
	}
}
