package main

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/abdheshnayak/gohtmlx/pkg/transpiler"
	"github.com/abdheshnayak/gohtmlx/pkg/utils"
)

// Set GOHTMLX_UPDATE_GOLDEN=1 to update golden files when the transpiler output intentionally changes.
const updateGoldenEnv = "GOHTMLX_UPDATE_GOLDEN"

func TestRun_DeterministicOutput(t *testing.T) {
	src := filepath.Join("example", "src", "comps")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("example source not found: %v", err)
	}

	dist1 := t.TempDir()
	dist2 := t.TempDir()

	utils.Log = utils.NewSlogLogger(slog.Default())
	opts := &transpiler.RunOptions{SingleFile: true}
	if err := transpiler.Run(src, dist1, opts); err != nil {
		t.Fatalf("first Run: %v", err)
	}
	if err := transpiler.Run(src, dist2, opts); err != nil {
		t.Fatalf("second Run: %v", err)
	}

	outFile := filepath.Join("gohtmlxc", "comp_generated.go")
	b1, err := os.ReadFile(filepath.Join(dist1, outFile))
	if err != nil {
		t.Fatalf("read first output: %v", err)
	}
	b2, err := os.ReadFile(filepath.Join(dist2, outFile))
	if err != nil {
		t.Fatalf("read second output: %v", err)
	}

	if len(b1) != len(b2) {
		t.Errorf("output length differs: first=%d second=%d", len(b1), len(b2))
	}
	for i := 0; i < len(b1) && i < len(b2); i++ {
		if b1[i] != b2[i] {
			t.Errorf("output differs at byte %d: %q vs %q", i, b1[i], b2[i])
			// Show context
			start := i - 20
			if start < 0 {
				start = 0
			}
			end := i + 20
			if end > len(b1) {
				end = len(b1)
			}
			t.Errorf("context first: %q", b1[start:end])
			t.Errorf("context second: %q", b2[start:end])
			break
		}
	}
}

func TestRun_SourceTrackingError(t *testing.T) {
	src := filepath.Join("testdata", "badprops")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("testdata not found: %v", err)
	}
	utils.Log = utils.NewSlogLogger(slog.Default())
	dist := t.TempDir()
	err := transpiler.Run(src, dist, &transpiler.RunOptions{SingleFile: true})
	if err == nil {
		t.Fatal("expected error from invalid props YAML")
	}
	var te *transpiler.TranspileError
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

func TestRun_Golden(t *testing.T) {
	src := filepath.Join("testdata", "golden")
	wantPath := filepath.Join("testdata", "golden", "want", "comp_generated.go")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("testdata/golden not found: %v", err)
	}
	if _, err := os.Stat(wantPath); err != nil {
		t.Skipf("golden want file not found: %v", err)
	}

	utils.Log = utils.NewSlogLogger(slog.Default())
	dist := t.TempDir()
	if err := transpiler.Run(src, dist, &transpiler.RunOptions{SingleFile: true}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	gotPath := filepath.Join(dist, "gohtmlxc", "comp_generated.go")
	got, err := os.ReadFile(gotPath)
	if err != nil {
		t.Fatalf("read generated: %v", err)
	}
	want, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}

	if string(got) != string(want) {
		if os.Getenv(updateGoldenEnv) == "1" {
			if err := os.WriteFile(wantPath, got, 0644); err != nil {
				t.Fatalf("update golden: %v", err)
			}
			t.Log("updated golden file")
			return
		}
		t.Errorf("generated output differs from golden. Run with %s=1 to update golden.", updateGoldenEnv)
		t.Logf("got %d bytes, want %d bytes", len(got), len(want))
		// Show first diff
		for i := 0; i < len(got) && i < len(want); i++ {
			if got[i] != want[i] {
				t.Errorf("first diff at byte %d", i)
				start := i - 30
				if start < 0 {
					start = 0
				}
				end := i + 30
				if end > len(got) {
					end = len(got)
				}
				t.Errorf("got  ...%q...", got[start:end])
				if end > len(want) {
					end = len(want)
				}
				t.Errorf("want ...%q...", want[start:end])
				break
			}
		}
	}
}

func TestRun_GoldenPerComponent(t *testing.T) {
	src := filepath.Join("testdata", "golden")
	wantDir := filepath.Join("testdata", "golden", "want_per_component")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("testdata/golden not found: %v", err)
	}
	if _, err := os.Stat(wantDir); err != nil {
		t.Skipf("golden want_per_component dir not found: %v", err)
	}

	utils.Log = utils.NewSlogLogger(slog.Default())
	dist := t.TempDir()
	if err := transpiler.Run(src, dist, &transpiler.RunOptions{SingleFile: false}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	gotDir := filepath.Join(dist, "gohtmlxc")

	entries, err := os.ReadDir(gotDir)
	if err != nil {
		t.Fatalf("read generated dir: %v", err)
	}
	var gotFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			gotFiles = append(gotFiles, e.Name())
		}
	}
	sort.Strings(gotFiles)

	for _, name := range gotFiles {
		gotPath := filepath.Join(gotDir, name)
		wantPath := filepath.Join(wantDir, name)
		got, err := os.ReadFile(gotPath)
		if err != nil {
			t.Fatalf("read generated %s: %v", name, err)
		}
		want, err := os.ReadFile(wantPath)
		if err != nil {
			if os.Getenv(updateGoldenEnv) == "1" {
				if err := os.WriteFile(wantPath, got, 0644); err != nil {
					t.Fatalf("update golden %s: %v", name, err)
				}
				t.Logf("updated golden %s", name)
				continue
			}
			t.Fatalf("read golden %s: %v", name, err)
		}
		if string(got) != string(want) {
			if os.Getenv(updateGoldenEnv) == "1" {
				if err := os.WriteFile(wantPath, got, 0644); err != nil {
					t.Fatalf("update golden %s: %v", name, err)
				}
				t.Logf("updated golden %s", name)
				return
			}
			t.Errorf("generated %s differs from golden. Run with %s=1 to update golden.", name, updateGoldenEnv)
			t.Logf("got %d bytes, want %d bytes", len(got), len(want))
		}
	}
}

func TestRun_IntegrationBuild(t *testing.T) {
	src := filepath.Join("testdata", "golden")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("testdata/golden not found: %v", err)
	}
	root := repoRoot()
	if root == "" {
		t.Skip("repo root (go.mod) not found")
	}
	// Transpile into a path under the repo so "go build" can resolve the module
	dist := filepath.Join(root, "testdata", "golden", "out")
	_ = os.RemoveAll(dist)
	defer os.RemoveAll(dist)
	utils.Log = utils.NewSlogLogger(slog.Default())
	if err := transpiler.Run(src, dist, &transpiler.RunOptions{SingleFile: true}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	relPkg, _ := filepath.Rel(root, filepath.Join(dist, "gohtmlxc"))
	pkgPath := "./" + filepath.ToSlash(relPkg)
	cmd := exec.Command("go", "build", pkgPath)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("go build %s: %v", pkgPath, err)
	}
}

func repoRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
