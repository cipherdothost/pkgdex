// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

// store_test.go
package metric_test

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"go.cipher.host/pkgdex/internal/database"
	"go.cipher.host/pkgdex/internal/metric"
)

func TestStore_Track(t *testing.T) { //nolint:paralleltest // we're tracking the number of installs
	var (
		tmpDir = t.TempDir()
		path   = filepath.Join(tmpDir, "track_test.db")
	)

	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	t.Cleanup(func() {
		if err = db.Close(); err != nil {
			t.Errorf("failed to close store: %v", err)
		}
	})

	store := metric.NewStore(db)

	tests := []struct {
		name           string
		givePkg        string
		giveTimestamp  time.Time
		wantErr        error
		wantCheckCount bool
		wantCount      int64
	}{
		{
			name:           "empty package returns error",
			givePkg:        "",
			giveTimestamp:  time.Now(),
			wantErr:        metric.ErrEmptyPackage,
			wantCheckCount: false,
		},
		{
			name:           "valid package succeeds",
			givePkg:        "example.com/user/package",
			giveTimestamp:  time.Now(),
			wantErr:        nil,
			wantCheckCount: true,
			wantCount:      1,
		},
		{
			name:    "duplicate tracking succeeds",
			givePkg: "example.com/user/package",
			giveTimestamp: func() time.Time {
				now := time.Now()

				return now.Add(1 + time.Millisecond)
			}(),
			wantErr:        nil,
			wantCheckCount: true,
			wantCount:      2,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // we're tracking the number of installs
		t.Run(tt.name, func(t *testing.T) {
			err = store.Track(tt.givePkg, tt.giveTimestamp)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Track() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("Track() unexpected error: %v", err)
			}

			if tt.wantCheckCount {
				var count int64

				count, err = store.Count(tt.givePkg)
				if err != nil {
					t.Errorf("Count() failed: %v", err)
				}

				if count != tt.wantCount {
					t.Errorf("Count() = %d, want %d", count, tt.wantCount)
				}
			}
		})
	}
}

func TestStore_Track_Concurrent(t *testing.T) {
	t.Parallel()

	var (
		tmpDir = t.TempDir()
		path   = filepath.Join(tmpDir, "concurrent_test.db")
	)

	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	t.Cleanup(func() {
		if err = db.Close(); err != nil {
			t.Errorf("failed to close store: %v", err)
		}
	})

	store := metric.NewStore(db)

	const (
		numGoroutines = 100
		pkg           = "example.com/user/package"
	)

	var (
		errChan = make(chan error, numGoroutines)
		workers = make([]func(), numGoroutines)
		wg      sync.WaitGroup
	)

	wg.Add(numGoroutines)

	for i := range workers {
		workers[i] = func() {
			defer wg.Done()

			// Add a small random delay to increase chance of concurrent access.
			time.Sleep(time.Duration(i) * time.Millisecond)

			if werr := store.Track(pkg, time.Now()); werr != nil {
				errChan <- fmt.Errorf("track failed: %w", werr)

				return
			}
		}
	}

	for _, worker := range workers {
		go worker()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("goroutine error: %v", err)
	}

	count, err := store.Count(pkg)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	// Due to the UNIQUE constraint and our handling of it, we expect the count
	// to be less than or equal to numGoroutines, but greater than 0.
	if count <= 0 || count > numGoroutines {
		t.Errorf("unexpected count = %d, want > 0 and <= %d", count, numGoroutines)
	}

	t.Logf("Final installation count: %d out of %d attempts", count, numGoroutines)
}

func TestStore_Count(t *testing.T) { //nolint:paralleltest // this leads to a race condition
	var (
		tmpDir = t.TempDir()
		path   = filepath.Join(tmpDir, "count_test.db")
	)

	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	t.Cleanup(func() {
		if err = db.Close(); err != nil {
			t.Errorf("failed to close store: %v", err)
		}
	})

	store := metric.NewStore(db)

	pkg := "example.com/user/package"
	for range 3 {
		if err = store.Track(pkg, time.Now()); err != nil {
			t.Fatalf("failed to setup test data: %v", err)
		}
	}

	tests := []struct {
		name    string
		givePkg string
		want    int64
		wantErr error
	}{
		{
			name:    "empty package returns error",
			givePkg: "",
			want:    0,
			wantErr: metric.ErrEmptyPackage,
		},
		{
			name:    "non-existent package returns zero",
			givePkg: "example.com/user/nonexistent",
			want:    0,
			wantErr: nil,
		},
		{
			name:    "existing package returns count",
			givePkg: pkg,
			want:    3,
			wantErr: nil,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // this leads to a race condition
		t.Run(tt.name, func(t *testing.T) {
			var got int64

			got, err = store.Count(tt.givePkg)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("Count() unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("Count() = %d, want %d", got, tt.want)
			}
		})
	}
}
