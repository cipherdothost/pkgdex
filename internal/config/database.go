// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/meta"
)

const (
	// ErrDatabaseMissingPath is returned when the database path is missing from
	// the configuration.
	ErrDatabaseMissingPath cmdkit.Error = "database path is missing"

	// ErrCreateDatabase is returned when the database file or its parent directories
	// cannot be created.
	ErrCreateDatabase cmdkit.Error = "could not create database"
)

const (
	defaultDatabasePath      = "/var/lib/" + meta.Name + "/" + meta.Name + ".db"
	defaultDatabaseIndexPath = "/var/lib/" + meta.Name + "/index"
)

// Database represents the database configuration.
type Database struct {
	// Path is the path to the database file.
	Path string `json:"path"`

	// IndexPath is the path to the search index database file.
	IndexPath string `json:"indexPath"`
}

// CreateIfNotExist creates the database directory and file if they don't exist.
// The directory is created with 0755 permissions and the file with 0600 permissions.
func (d *Database) CreateIfNotExist() error {
	if d.Path == "" {
		return ErrDatabaseMissingPath
	}

	dir := filepath.Dir(d.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("%w: failed to create directory %s: %w", ErrCreateDatabase, dir, err)
	}

	f, err := os.OpenFile(d.Path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("%w: failed to create database file: %w", ErrCreateDatabase, err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("%w: failed to close database file: %w", ErrCreateDatabase, err)
	}

	return nil
}

// validate validates the database configuration.
func (d *Database) validate() error {
	if d.Path == "" {
		return ErrDatabaseMissingPath
	}

	return nil
}

// setDefaults sets the default values for the database configuration.
func (d *Database) setDefaults() {
	if d.Path == "" {
		d.Path = defaultDatabasePath
	}

	if d.IndexPath == "" {
		d.IndexPath = defaultDatabaseIndexPath
	}
}
