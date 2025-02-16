// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"go.cipher.host/pkgdex/internal/config"
)

func TestDatabase_CreateIfNotExist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(t *testing.T) (*config.Database, func())
		wantErr   error
		wantPerms fs.FileMode
	}{
		{
			name: "success - creates new database file and directory",
			setup: func(t *testing.T) (*config.Database, func()) {
				t.Helper()

				var (
					tmpDir = t.TempDir()
					dbPath = filepath.Join(tmpDir, "subdir", "test.db")
				)

				return &config.Database{
					Path: dbPath,
				}, func() {}
			},
			wantErr:   nil,
			wantPerms: 0o600,
		},
		{
			name: "success - file already exists",
			setup: func(t *testing.T) (*config.Database, func()) {
				t.Helper()

				var (
					tmpDir = t.TempDir()
					dbPath = filepath.Join(tmpDir, "test.db")
				)

				if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}

				if err := os.WriteFile(dbPath, make([]byte, 0), 0o600); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}

				return &config.Database{
					Path: dbPath,
				}, func() {}
			},
			wantErr:   nil,
			wantPerms: 0o600,
		},
		{
			name: "error - empty path",
			setup: func(t *testing.T) (*config.Database, func()) {
				t.Helper()

				return &config.Database{
					Path: "",
				}, func() {}
			},
			wantErr: config.ErrDatabaseMissingPath,
		},
		{
			name: "error - permission denied for directory",
			setup: func(t *testing.T) (*config.Database, func()) {
				t.Helper()

				var (
					tmpDir      = t.TempDir()
					readOnlyDir = filepath.Join(tmpDir, "readonly")
				)

				if err := os.Mkdir(readOnlyDir, 0o555); err != nil {
					t.Fatalf("failed to create read-only directory: %v", err)
				}

				dbPath := filepath.Join(readOnlyDir, "subdir", "test.db")

				return &config.Database{
						Path: dbPath,
					}, func() {
						os.Chmod(readOnlyDir, 0o755) //nolint:errcheck // cleanup
					}
			},
			wantErr: config.ErrCreateDatabase,
		},
		{
			name: "error - permission denied for file",
			setup: func(t *testing.T) (*config.Database, func()) {
				t.Helper()

				tmpDir := t.TempDir()

				if err := os.Chmod(tmpDir, 0o555); err != nil {
					t.Fatalf("failed to set directory permissions: %v", err)
				}

				dbPath := filepath.Join(tmpDir, "test.db")

				return &config.Database{
						Path: dbPath,
					}, func() {
						os.Chmod(tmpDir, 0o755) //nolint:errcheck // cleanup
					}
			},
			wantErr: config.ErrCreateDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, cleanup := tt.setup(t)
			defer cleanup()

			err := db.CreateIfNotExist()
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("CreateIfNotExist() error = %v, wantErr %v", err, tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("CreateIfNotExist() unexpected error: %v", err)
			}

			info, err := os.Stat(db.Path)
			if err != nil {
				t.Fatalf("failed to stat database file: %v", err)
			}

			if info.Mode().Perm() != tt.wantPerms {
				t.Errorf("database file has wrong permissions: got %v, want %v",
					info.Mode().Perm(), tt.wantPerms)
			}

			dirInfo, err := os.Stat(filepath.Dir(db.Path))
			if err != nil {
				t.Fatalf("failed to stat database directory: %v", err)
			}

			if !dirInfo.IsDir() {
				t.Error("database parent path is not a directory")
			}

			if dirInfo.Mode().Perm() != 0o755 {
				t.Errorf("database directory has wrong permissions: got %v, want %v",
					dirInfo.Mode().Perm(), 0o755)
			}
		})
	}
}
