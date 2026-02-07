package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// HTMLFile holds a single .html file path and its content (for source tracking).
type HTMLFile struct {
	Path    string
	Content []byte
}

// ReadFileToByteArray reads the file at the given path and returns its content as a byte array
func ReadFileToByteArray(path string) ([]byte, error) {
	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", path, err)
	}
	return content, nil
}

// WalkAndReadHTMLFiles walks the directory recursively and returns each .html file's path and content.
// Files are returned sorted by path for deterministic behavior.
func WalkAndReadHTMLFiles(dir string) ([]HTMLFile, error) {
	var files []HTMLFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error walking the path", path, ":", err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
			content, err := ReadFileToByteArray(path)
			if err != nil {
				log.Println("Error reading file", path, ":", err)
				return err
			}
			files = append(files, HTMLFile{Path: path, Content: content})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the directory: %v", err)
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

// WalkAndConcatenateHTML walks through the directory recursively and concatenates all .html file contents.
// Deprecated: use WalkAndReadHTMLFiles for source tracking; this is kept for compatibility.
func WalkAndConcatenateHTML(dir string) ([]byte, error) {
	files, err := WalkAndReadHTMLFiles(dir)
	if err != nil {
		return nil, err
	}
	var out []byte
	for _, f := range files {
		out = append(out, f.Content...)
		out = append(out, '\n')
	}
	return out, nil
}
