package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetPaths(dir string, ext string) ([]string, error) {
	paths := []string{}

	// Walk the directory recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error walking the path", path, ":", err)
			return err
		}

		// Check if file has .html extension
		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the directory: %v", err)
	}

	return paths, nil
}
