package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ReadFileToByteArray reads the file at the given path and returns its content as a byte array
func ReadFileToByteArray(path string) ([]byte, error) {
	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", path, err)
	}
	return content, nil
}

// WalkAndConcatenateHTML walks through the directory recursively and concatenates all .html file contents
func WalkAndConcatenateHTML(dir string) ([]byte, error) {
	var concatenatedContent []byte

	// Walk the directory recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error walking the path", path, ":", err)
			return err
		}

		// Check if file has .html extension
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
			// Read file content and append it to the concatenatedContent
			content, err := ReadFileToByteArray(path)
			if err != nil {
				log.Println("Error reading file", path, ":", err)
				return err
			}
			concatenatedContent = append(concatenatedContent, content...)
			concatenatedContent = append(concatenatedContent, '\n') // Add newline between files
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the directory: %v", err)
	}

	return concatenatedContent, nil
}
