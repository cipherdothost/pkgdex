// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package metric

import (
	"bytes"
	"fmt"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/database"
	bolt "go.etcd.io/bbolt"
)

const (
	// ErrEmptyPackage is returned when an empty package name is provided.
	ErrEmptyPackage cmdkit.Error = "package name cannot be empty; please provide a valid package name"

	// ErrTrackInstall is returned when an installation cannot be tracked.
	ErrTrackInstall cmdkit.Error = "failed to record package installation; database may be unavailable"

	// ErrGetCount is returned when the installation count cannot be retrieved.
	ErrGetCount cmdkit.Error = "failed to retrieve installation count for package"

	// ErrBucketNotFound is returned when a database bucket is not found.
	ErrBucketNotFound cmdkit.Error = "database bucket not found; database may need reinitialization"
)

// Store represents a package installation metric store.
type Store struct {
	db *database.Store
}

// NewStore creates a new Store with the given database.
func NewStore(db *database.Store) *Store {
	return &Store{
		db: db,
	}
}

// Track records a new package installation in the database.
func (s *Store) Track(pkg string, timestamp time.Time) error {
	if pkg == "" {
		return ErrEmptyPackage
	}

	err := s.db.DB().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(xunsafe.StringToBytes(database.BucketInstalls))
		if b == nil {
			return fmt.Errorf("%w", ErrBucketNotFound)
		}

		key := fmt.Sprintf("%s:%d", pkg, timestamp.UTC().UnixNano())

		// We store an empty value because we only care about the count.
		if err := b.Put(xunsafe.StringToBytes(key), make([]byte, 0)); err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTrackInstall, err)
	}

	return nil
}

// Count returns the total number of installations for a package.
func (s *Store) Count(pkg string) (int64, error) {
	if pkg == "" {
		return 0, ErrEmptyPackage
	}

	var count int64

	err := s.db.DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(xunsafe.StringToBytes(database.BucketInstalls))
		if bucket == nil {
			return fmt.Errorf("%w", ErrBucketNotFound)
		}

		var (
			prefix = xunsafe.StringToBytes(pkg + ":")
			cursor = bucket.Cursor()
		)

		for key, _ := cursor.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, _ = cursor.Next() {
			count++
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("%w %q: %w", ErrGetCount, pkg, err)
	}

	return count, nil
}
