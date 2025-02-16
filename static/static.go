// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

// Package static provides static files and HTML templates for the service.
package static

import (
	"crypto/sha512"
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"path"
)

//go:embed *
var files embed.FS

// Files returns the embedded files.
func Files() embed.FS {
	return files
}

// Parse parses the templates and returns a template.Template instance.
func Parse() (*template.Template, error) {
	paths, err := GetStylePaths()
	if err != nil {
		return nil, fmt.Errorf("getting style paths: %w", err)
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mainStylePath": func() string {
			return paths["main.css"]
		},
		"highlightStylePath": func() string {
			return paths["highlight.css"]
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(Files(), "templates/*.html", "templates/components/*.html")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return tmpl, nil
}

// VersionHash generates a short hash of a file's contents to use for
// versioning.
func VersionHash(filepath string) (string, error) {
	content, err := files.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	hash := sha512.Sum512_256(content)

	return hex.EncodeToString(hash[:8]), nil
}

// VersionedPath returns the path to a file with its content hash included.
func VersionedPath(filepath string) (string, error) {
	version, err := VersionHash(filepath)
	if err != nil {
		return "", err
	}

	var (
		extension = path.Ext(filepath)
		base      = filepath[:len(filepath)-len(extension)]
	)

	return base + "." + version + extension, nil
}

// GetStylePaths returns the versioned paths for CSS files.
func GetStylePaths() (map[string]string, error) {
	files := []string{
		"assets/css/main.css",
		"assets/css/highlight.css",
	}

	paths := make(map[string]string, len(files))

	for _, file := range files {
		vpath, err := VersionedPath(file)
		if err != nil {
			return nil, fmt.Errorf("generating versioned path for %q: %w", file, err)
		}

		paths[path.Base(file)] = vpath
	}

	return paths, nil
}
