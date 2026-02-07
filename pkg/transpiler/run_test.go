package transpiler

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// testdata paths relative to repo root; tests may skip if not found
const (
	goldenSrc  = "testdata/golden"
	badpropsSrc = "testdata/badprops"
)

func findTestdata(t *testing.T, subpath string) string {
	t.Helper()
	// Run from repo root or from pkg/transpiler
	for _, base := range []string{"../../", "./"} {
		dir := filepath.Join(base, subpath)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	t.Skipf("testdata not found: %s", subpath)
	return ""
}

func TestRun_SingleFile_true(t *testing.T) {
	src := findTestdata(t, goldenSrc)
	dist := t.TempDir()

	err := Run(src, dist, &RunOptions{SingleFile: true, Pkg: "gohtmlxc"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	outPath := filepath.Join(dist, "gohtmlxc", "comp_generated.go")
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(content), "package gohtmlxc") {
		t.Errorf("expected package gohtmlxc in output, got:\n%s", content)
	}
	if !strings.Contains(string(content), "func AComp(") || !strings.Contains(string(content), "func BComp(") {
		t.Errorf("expected AComp and BComp in single file, got:\n%s", content)
	}
}

func TestRun_SingleFile_false(t *testing.T) {
	src := findTestdata(t, goldenSrc)
	dist := t.TempDir()

	err := Run(src, dist, &RunOptions{SingleFile: false, Pkg: "gohtmlxc"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	outDir := filepath.Join(dist, "gohtmlxc")
	aPath := filepath.Join(outDir, "A.go")
	bPath := filepath.Join(outDir, "B.go")
	for _, p := range []string{aPath, bPath} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected file %s to exist: %v", p, err)
		}
	}
	aContent, _ := os.ReadFile(aPath)
	bContent, _ := os.ReadFile(bPath)
	if !strings.Contains(string(aContent), "package gohtmlxc") {
		t.Errorf("A.go should contain package gohtmlxc")
	}
	if !strings.Contains(string(bContent), "package gohtmlxc") {
		t.Errorf("B.go should contain package gohtmlxc")
	}
	if !strings.Contains(string(aContent), "func AComp(") {
		t.Errorf("A.go should contain AComp")
	}
	if !strings.Contains(string(bContent), "func BComp(") {
		t.Errorf("B.go should contain BComp")
	}
}

func TestRun_CustomPkg(t *testing.T) {
	src := findTestdata(t, goldenSrc)
	dist := t.TempDir()
	customPkg := "mypkg"

	err := Run(src, dist, &RunOptions{SingleFile: true, Pkg: customPkg})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	outPath := filepath.Join(dist, customPkg, "comp_generated.go")
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(content), "package "+customPkg) {
		t.Errorf("expected package %s in output, got:\n%s", customPkg, content)
	}
}

func TestRun_ErrorReturnsTranspileError(t *testing.T) {
	src := findTestdata(t, badpropsSrc)
	dist := t.TempDir()

	err := Run(src, dist, &RunOptions{SingleFile: true})
	if err == nil {
		t.Fatal("expected error from invalid props YAML")
	}

	var te *TranspileError
	if !errors.As(err, &te) {
		t.Fatalf("expected TranspileError, got %T: %v", err, err)
	}
	if te.FilePath == "" {
		t.Error("TranspileError should have FilePath set")
	}
	if !strings.Contains(te.FilePath, "bad.html") {
		t.Errorf("FilePath should mention bad.html, got %q", te.FilePath)
	}
	if te.Message == "" {
		t.Error("TranspileError should have Message set")
	}
}

func TestRun_NilOptionsUsesDefaults(t *testing.T) {
	src := findTestdata(t, goldenSrc)
	dist := t.TempDir()

	err := Run(src, dist, nil)
	if err != nil {
		t.Fatalf("Run with nil options: %v", err)
	}

	// Default: one file per component, pkg "gohtmlxc"
	outDir := filepath.Join(dist, "gohtmlxc")
	if _, err := os.Stat(outDir); err != nil {
		t.Fatalf("expected dist/gohtmlxc: %v", err)
	}
	aPath := filepath.Join(outDir, "A.go")
	if _, err := os.Stat(aPath); err != nil {
		t.Errorf("expected A.go with nil options: %v", err)
	}
}

func TestRun_IncrementalSkipsWhenOutputNewer(t *testing.T) {
	src := findTestdata(t, goldenSrc)
	dist := t.TempDir()

	// First run: populate dist (generated .go files get current mtime)
	if err := Run(src, dist, &RunOptions{SingleFile: false}); err != nil {
		t.Fatalf("first Run: %v", err)
	}
	// Second run with Incremental: no .html changed, so should skip (return nil)
	if err := Run(src, dist, &RunOptions{Incremental: true}); err != nil {
		t.Fatalf("incremental Run should skip and return nil: %v", err)
	}
}
