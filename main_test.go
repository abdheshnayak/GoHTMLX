package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_DeterministicOutput(t *testing.T) {
	src := filepath.Join("example", "src", "comps")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("example source not found: %v", err)
	}

	dist1 := t.TempDir()
	dist2 := t.TempDir()

	if err := Run(src, dist1); err != nil {
		t.Fatalf("first Run: %v", err)
	}
	if err := Run(src, dist2); err != nil {
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
