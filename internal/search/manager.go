// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package search

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2"
	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/config"
)

const (
	// ErrOpenIndex is returned when the search index database cannot be opened.
	ErrOpenIndex cmdkit.Error = "failed to open search index; check file permissions and disk space"

	// ErrCloseIndex is returned when the search index database cannot be closed.
	ErrCloseIndex cmdkit.Error = "failed to close search index; some operations may not have been saved"

	// ErrCreateIndex is returned when the search index database cannot be created.
	ErrCreateIndex cmdkit.Error = "failed to create search index; check permissions and available space"

	// ErrIndexPackage is returned when a package cannot be indexed.
	ErrIndexPackage cmdkit.Error = "failed to add package to search index"

	// ErrSearchPackages is returned when packages cannot be searched.
	ErrSearchPackages cmdkit.Error = "search operation failed; the search index may be corrupted or unavailable"
)

// Package represents a package in the search index.
type Package struct {
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	License     string    `json:"license"`
	Repository  string    `json:"repository"`
}

// Manager represents a search index manager.
type Manager struct {
	index bleve.Index
	mu    sync.RWMutex
}

// NewManager creates a new search index Manager instance.
func NewManager(path string) (*Manager, error) {
	index, err := bleve.Open(path)
	if err != nil {
		if errors.Is(err, bleve.ErrorIndexPathDoesNotExist) {
			var (
				mapping    = bleve.NewIndexMapping()
				docMapping = bleve.NewDocumentMapping()
			)

			nameField := bleve.NewTextFieldMapping()
			nameField.Store = true
			nameField.Index = true
			nameField.IncludeTermVectors = true

			docMapping.AddFieldMappingsAt("name", nameField)

			descField := bleve.NewTextFieldMapping()
			descField.Store = true
			descField.Index = true
			descField.IncludeTermVectors = true

			docMapping.AddFieldMappingsAt("description", descField)

			mapping.AddDocumentMapping("package", docMapping)

			index, err = bleve.New(path, mapping)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrCreateIndex, err)
			}

			return &Manager{
				index: index,
			}, nil
		}

		return nil, fmt.Errorf("%w: %w", ErrOpenIndex, err)
	}

	return &Manager{
		index: index,
	}, nil
}

// Close closes the search index.
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.index.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrCloseIndex, err)
	}

	return nil
}

// Index adds or updates a package in the search index.
func (m *Manager) Index(pkg *config.Package, updatedAt time.Time) error {
	if pkg.Hidden {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	doc := Package{
		Name:        pkg.Name,
		Description: pkg.Description,
		License:     pkg.License,
		Repository:  pkg.Repository,
		UpdatedAt:   updatedAt,
	}

	if err := m.index.Index(pkg.Name, doc); err != nil {
		return fmt.Errorf("%w (%q): %w", ErrIndexPackage, pkg.Name, err)
	}

	return nil
}

// Search performs a search query and returns matching packages.
func (m *Manager) Search(query string) ([]Package, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	q := bleve.NewDisjunctionQuery(
		bleve.NewMatchPhraseQuery(query),
		bleve.NewFuzzyQuery(query),
	)

	search := bleve.NewSearchRequest(q)
	search.Fields = []string{
		"name",
		"description",
		"license",
		"repository",
		"updated_at",
	}
	search.Highlight = bleve.NewHighlight()

	result, err := m.index.Search(search)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSearchPackages, err)
	}

	docs := make([]Package, 0, len(result.Hits))

	for _, hit := range result.Hits {
		var doc Package

		if name, ok := hit.Fields["name"].(string); ok {
			doc.Name = name
		}

		if desc, ok := hit.Fields["description"].(string); ok {
			doc.Description = desc
		}

		if license, ok := hit.Fields["license"].(string); ok {
			doc.License = license
		}

		if repo, ok := hit.Fields["repository"].(string); ok {
			doc.Repository = repo
		}

		if updated, ok := hit.Fields["updated_at"].(time.Time); ok {
			doc.UpdatedAt = updated
		}

		docs = append(docs, doc)
	}

	return docs, nil
}
