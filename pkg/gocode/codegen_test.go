package gocode

import (
	"strings"
	"testing"
)

func TestConstructStruct(t *testing.T) {
	props := map[string]string{"name": "string", "count": "int"}
	out := ConstructStruct(props, "Widget")
	if !strings.Contains(out, "type Widget struct") {
		t.Errorf("output should contain type Widget struct, got:\n%s", out)
	}
	if !strings.Contains(out, "Name") || !strings.Contains(out, "string") {
		t.Errorf("output should contain Name and string (capitalized), got:\n%s", out)
	}
	if !strings.Contains(out, "Count") || !strings.Contains(out, "int") {
		t.Errorf("output should contain Count and int, got:\n%s", out)
	}
	if !strings.Contains(out, "Attrs Attrs") {
		t.Errorf("output should contain Attrs Attrs, got:\n%s", out)
	}
	// Deterministic: same props should produce same output on repeated calls
	out2 := ConstructStruct(props, "Widget")
	if out != out2 {
		t.Error("ConstructStruct should be deterministic")
	}
}

func TestConstructStruct_EmptyProps(t *testing.T) {
	out := ConstructStruct(map[string]string{}, "Empty")
	if !strings.Contains(out, "type Empty struct") {
		t.Errorf("output should contain type Empty struct, got:\n%s", out)
	}
	if !strings.Contains(out, "Attrs Attrs") {
		t.Errorf("output should contain Attrs Attrs, got:\n%s", out)
	}
}

func TestConstructSource(t *testing.T) {
	// Use minimal valid code; struct must end with newline like ConstructStruct output
	codes := map[string]string{"Foo": "R(E(`div`, Attrs{},))"}
	structs := []string{"type Foo struct {\n\tAttrs Attrs\n}\n"}
	imports := []string{}
	out, err := ConstructSource(codes, structs, imports)
	if err != nil {
		t.Fatalf("ConstructSource: %v", err)
	}
	if !strings.Contains(out, "package gohtmlxc") {
		t.Errorf("output should contain package gohtmlxc, got:\n%s", out)
	}
	if !strings.Contains(out, "func FooComp(") {
		t.Errorf("output should contain FooComp, got:\n%s", out)
	}
	// Deterministic
	out2, err := ConstructSource(codes, structs, imports)
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if out != out2 {
		t.Error("ConstructSource should be deterministic")
	}
}
