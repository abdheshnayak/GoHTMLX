package utils

import (
	"strings"
	"testing"
	"text/template"
)

func TestParseSections(t *testing.T) {
	// Same delimiters as main.go for component sections
	tmpl, err := template.New("root").Delims("<!-- +", " -->").Parse(
		`<!-- + define "A" -->content A<!-- + end -->` +
			`<!-- + define "B" -->content B<!-- + end -->`,
	)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	sections, err := ParseSections(tmpl)
	if err != nil {
		t.Fatalf("ParseSections: %v", err)
	}
	if len(sections) == 0 {
		t.Fatal("expected at least one section")
	}
	// Template names include the define directive; we only care that content is extracted
	hasA := false
	for name, content := range sections {
		if strings.Contains(name, "A") && strings.TrimSpace(content) == "content A" {
			hasA = true
			break
		}
	}
	if !hasA {
		t.Errorf("expected section with content A; got sections: %v", sections)
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"foo", "Foo"},
		{"id", "Id"},
		{"", ""},
		{"a", "A"},
	}
	for _, tt := range tests {
		got := Capitalize(tt.in)
		if got != tt.want {
			t.Errorf("Capitalize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
