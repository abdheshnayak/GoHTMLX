package validate

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.html")
	content := `<!-- + define "X" -->
<div>ok</div>
<!-- + end -->
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Run(dir, &buf); err != nil {
		t.Fatalf("expected no error for valid file: %v\n%s", err, buf.Bytes())
	}
}

func TestRun_UnclosedBlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.html")
	content := `<!-- + define "X" -->
<div>no end</div>
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Run(dir, &buf); err == nil {
		t.Fatal("expected error for unclosed block")
	}
	out := buf.String()
	if out == "" {
		t.Error("expected validation output to stderr")
	}
	if !bytes.Contains(buf.Bytes(), []byte("unclosed")) {
		t.Errorf("output should mention unclosed: %q", out)
	}
}

func TestRun_MismatchedEnd(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mismatch.html")
	content := `<!-- + define "X" -->
<div>ok</div>
<!-- | end -->
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Run(dir, &buf); err == nil {
		t.Fatal("expected error for mismatched end")
	}
	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("does not match")) {
		t.Errorf("output should mention mismatch: %q", out)
	}
}
