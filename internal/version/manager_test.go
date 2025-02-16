// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package version_test

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/database"
	"go.cipher.host/pkgdex/internal/version"
)

func TestManager_TrackPackage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		givePkg       *config.Package
		giveTimestamp time.Time
		setup         func(t *testing.T, manager *version.Manager) error
		wantErr       error
	}{
		{
			name: "valid package succeeds",
			givePkg: &config.Package{
				Name:    "example.com/user/package",
				Version: "1.0.0",
			},
			giveTimestamp: time.Now(),
			wantErr:       nil,
		},
		{
			name: "duplicate version is ignored",
			givePkg: &config.Package{
				Name:    "example.com/user/package",
				Version: "1.0.0",
			},
			giveTimestamp: time.Now(),
			setup: func(t *testing.T, manager *version.Manager) error {
				t.Helper()

				return manager.TrackPackage(&config.Package{
					Name:    "example.com/user/package",
					Version: "1.0.0",
				}, time.Now())
			},
			wantErr: nil,
		},
		{
			name: "new version is tracked",
			givePkg: &config.Package{
				Name:    "example.com/user/package",
				Version: "1.1.0",
			},
			giveTimestamp: time.Now(),
			setup: func(t *testing.T, manager *version.Manager) error {
				t.Helper()

				return manager.TrackPackage(&config.Package{
					Name:    "example.com/user/package",
					Version: "1.0.0",
				}, time.Now())
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				tmpDir = t.TempDir()
				path   = filepath.Join(tmpDir, "track_test.db")
			)

			db, err := database.Open(path)
			if err != nil {
				t.Fatalf("TrackPackage() failed to open database: %v", err)
			}

			t.Cleanup(func() {
				if err = db.Close(); err != nil {
					t.Errorf("TrackPackage() failed to close store: %v", err)
				}
			})

			manager := version.NewManager(db)

			if tt.setup != nil {
				if err = tt.setup(t, manager); err != nil {
					t.Fatalf("TrackPackage() setup failed: %v", err)
				}
			}

			err = manager.TrackPackage(tt.givePkg, tt.giveTimestamp)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("TrackPackage() error = %v, wantErr %v", err, tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Errorf("TrackPackage() unexpected error: %v", err)
			}
		})
	}
}

func TestManager_GetLastUpdate(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name          string
		givePkg       *config.Package
		giveTimestamp time.Time
		setup         func(t *testing.T, manager *version.Manager) error
		want          time.Time
	}{
		{
			name: "no history returns given timestamp",
			givePkg: &config.Package{
				Name:    "example.com/user/package",
				Version: "1.0.0",
			},
			giveTimestamp: now,
			want:          now,
		},
		{
			name: "returns latest update time",
			givePkg: &config.Package{
				Name:    "example.com/user/package",
				Version: "1.1.0",
			},
			giveTimestamp: now,
			setup: func(t *testing.T, manager *version.Manager) error {
				t.Helper()

				pkg := &config.Package{
					Name:    "example.com/user/package",
					Version: "1.0.0",
				}

				oldTime := now.Add(-24 * time.Hour)

				if err := manager.TrackPackage(pkg, oldTime); err != nil {
					return fmt.Errorf("failed to track old version: %w", err)
				}

				pkg.Version = "1.1.0"

				if err := manager.TrackPackage(pkg, now); err != nil {
					return fmt.Errorf("failed to track new version: %w", err)
				}

				return nil
			},
			want: now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				tmpDir = t.TempDir()
				path   = filepath.Join(tmpDir, "last_update_test.db")
			)

			db, err := database.Open(path)
			if err != nil {
				t.Fatalf("GetLastUpdate() failed to open database: %v", err)
			}

			t.Cleanup(func() {
				if err = db.Close(); err != nil {
					t.Errorf("GetLastUpdate() failed to close store: %v", err)
				}
			})

			manager := version.NewManager(db)

			if tt.setup != nil {
				if err = tt.setup(t, manager); err != nil {
					t.Fatalf("GetLastUpdate() test setup failed: %v", err)
				}
			}

			got := manager.GetLastUpdate(tt.givePkg, tt.giveTimestamp)
			if !got.Equal(tt.want) {
				t.Errorf("GetLastUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
