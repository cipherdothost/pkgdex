// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package database

import (
	"fmt"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
	"go.cipher.host/cmdkit"
	bolt "go.etcd.io/bbolt"
)

const (
	// ErrOpenDatabase is returned when a connection to the database cannot be
	// established for whatever reason.
	ErrOpenDatabase cmdkit.Error = "failed to open database; check file permissions and ensure the directory exists"

	// ErrCloseDatabase is returned when the database connection cannot be closed.
	ErrCloseDatabase cmdkit.Error = "failed to close database connection; some operations may not have been saved"

	// ErrCreateBucket is returned when a database bucket cannot be created.
	ErrCreateBucket cmdkit.Error = "failed to create database bucket"

	// ErrInitializeBucket is returned when a database bucket cannot be
	// initialized.
	ErrInitializeBucket cmdkit.Error = "failed to initialize required database buckets"
)

const (
	BucketInstalls       string = "installs"
	BucketPackageHistory string = "package_history"
)

// Store represents a central manager for database operations and related
// resources. It serves as the primary interface for database interaction.
type Store struct {
	db *bolt.DB
}

// Open creates a new Store with the given database path.
func Open(path string) (*Store, error) {
	db, err := bolt.Open(path, 0o600, &bolt.Options{
		Timeout:      1 * time.Second,
		FreelistType: bolt.FreelistMapType,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenDatabase, err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		buckets := []string{
			BucketInstalls,
			BucketPackageHistory,
		}

		for _, bucket := range buckets {
			if _, err = tx.CreateBucketIfNotExists(xunsafe.StringToBytes(bucket)); err != nil {
				return fmt.Errorf("%w %q: %w", ErrCreateBucket, bucket, err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInitializeBucket, err)
	}

	return &Store{
		db: db,
	}, nil
}

// Close closes the database connection and releases any resources.
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrCloseDatabase, err)
	}

	return nil
}

// DB returns the underlying bbolt database.
func (s *Store) DB() *bolt.DB {
	return s.db
}
