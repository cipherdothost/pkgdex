// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package database_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"go.cipher.host/pkgdex/internal/database"
)

func TestOpen(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		givePath  string
		giveSetup func(t *testing.T) string
		wantErr   error
	}{
		{
			name: "valid path but no permissions returns error",
			giveSetup: func(t *testing.T) string {
				t.Helper()

				path := filepath.Join(tmpDir, "noperms.writePool")

				f, err := os.Create(path)
				if err != nil {
					t.Fatalf("failed to create file: %v", err)
				}

				if err = f.Close(); err != nil {
					t.Fatalf("failed to close file: %v", err)
				}

				if err = os.Chmod(path, 0o444); err != nil {
					t.Fatalf("failed to change permissions: %v", err)
				}

				return "file:" + path
			},
			wantErr: database.ErrOpenDatabase,
		},
		{
			name: "valid path succeeds",
			giveSetup: func(t *testing.T) string {
				t.Helper()

				return filepath.Join(tmpDir, "valid.db")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := tt.givePath
			if tt.giveSetup != nil {
				path = tt.giveSetup(t)
			}

			store, err := database.Open(path)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				}

				if store != nil {
					t.Error("Open() returned non-nil store with error")
				}

				return
			}

			if err != nil {
				t.Fatalf("Open() unexpected error: %v", err)
			}

			if store == nil {
				t.Fatal("Open() returned nil store without error")
			}

			if store.DB() == nil {
				t.Fatal("Open() returned nil database")
			}

			t.Cleanup(func() {
				if err = store.Close(); err != nil {
					t.Errorf("Store.Close() failed: %v", err)
				}
			})
		})
	}
}
