package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func main() {
	// Detect the first chart folder under /templates
	dirs, err := os.ReadDir("/templates")
	if err != nil {
		panic(err)
	}

	var chart string
	for _, d := range dirs {
		if d.IsDir() {
			chart = d.Name()
			break
		}
	}
	if chart == "" {
		panic("no chart folder found in /templates")
	}

	root := filepath.Join("/templates", chart, "templates")

	images, err := ExtractImages(root)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found images:")
	for _, img := range images {
		fmt.Println(" -", img)
	}

	// Get output folder from environment variable
	reportsPath := os.Getenv("OUTPUT_FOLDER")
	if reportsPath == "" {
		reportsPath = "/reports" // fallback
	}

	// Ensure the folder exists
	if err := os.MkdirAll(reportsPath, 0755); err != nil {
		panic(fmt.Errorf("failed to create reports folder: %w", err))
	}

	imagesFile := filepath.Join(reportsPath, "images.txt")
	f, err := os.Create(imagesFile)
	if err != nil {
		panic(fmt.Errorf("failed to create images.txt: %w", err))
	}
	defer f.Close()

	for _, img := range images {
		if _, err := f.WriteString(img + "\n"); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write image %s: %v\n", img, err)
		}
	}
}

// ExtractImages walks a folder recursively and finds all container images
func ExtractImages(root string) ([]string, error) {
	var images []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		docs := splitYAMLDocuments(data)
		for _, doc := range docs {
			var m map[string]interface{}
			if err := yaml.Unmarshal(doc, &m); err != nil {
				continue
			}
			images = append(images, findImages(m)...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return unique(images), nil
}

// Split multi-doc YAML into individual docs
func splitYAMLDocuments(raw []byte) [][]byte {
	parts := bytes.Split(raw, []byte("\n---"))
	var docs [][]byte
	for _, p := range parts {
		cleaned := bytes.TrimSpace(p)
		if len(cleaned) > 0 {
			docs = append(docs, cleaned)
		}
	}
	return docs
}

// Recursively find all image strings in a YAML structure
func findImages(node interface{}) []string {
	var out []string
	switch n := node.(type) {
	case map[string]interface{}:
		for k, v := range n {
			if k == "image" {
				if s, ok := v.(string); ok {
					out = append(out, s)
				}
			}
			out = append(out, findImages(v)...)
		}
	case []interface{}:
		for _, v := range n {
			out = append(out, findImages(v)...)
		}
	}
	return out
}

// Return unique strings
func unique(items []string) []string {
	m := map[string]bool{}
	var out []string
	for _, i := range items {
		if !m[i] {
			m[i] = true
			out = append(out, i)
		}
	}
	return out
}

