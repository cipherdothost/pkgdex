// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package version

import (
	"encoding/json"
	"fmt"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
	jsoniter "github.com/json-iterator/go"
	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/database"
	bolt "go.etcd.io/bbolt"
)

const (
	// ErrTrackVersion is returned when a version cannot be tracked.
	ErrTrackVersion cmdkit.Error = "failed to track version for package"

	// ErrGetHistory is returned when the version history cannot be retrieved.
	ErrGetHistory cmdkit.Error = "failed to retrieve version history for package; check database connection"

	// ErrCreatePackageBucket is returned when a package database bucket cannot be created.
	ErrCreatePackageBucket cmdkit.Error = "failed to create package bucket"

	// ErrMarshalVersionData is returned when version data cannot be marshaled.
	ErrMarshalVersionData cmdkit.Error = "failed to process version data for package; data may be corrupted"

	// ErrStoreVersion is returned when a version cannot be stored in the
	// database.
	ErrStoreVersion cmdkit.Error = "failed to store package version in database"

	// ErrBucketNotFound is returned when a database bucket is not found.
	ErrBucketNotFound cmdkit.Error = "bucket not found"
)

// PackageVersion represents a version record in the database.
type PackageVersion struct {
	CreatedAt time.Time `json:"created_at"`
	Version   string    `json:"version"`
}

// Manager handles package version tracking.
type Manager struct {
	db *database.Store
}

// NewManager creates a new version manager.
func NewManager(db *database.Store) *Manager {
	return &Manager{
		db: db,
	}
}

// TrackPackage records a new package version if it has changed.
func (m *Manager) TrackPackage(pkg *config.Package, timestamp time.Time) error {
	err := m.db.DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(xunsafe.StringToBytes(database.BucketPackageHistory))
		if bucket == nil {
			return fmt.Errorf(" %q: %w", pkg.Name, ErrBucketNotFound)
		}

		pkgBucket, err := bucket.CreateBucketIfNotExists(xunsafe.StringToBytes(pkg.Name))
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCreatePackageBucket, err)
		}

		pkgVersion := xunsafe.StringToBytes(pkg.Version)

		if version := pkgBucket.Get(pkgVersion); version != nil {
			return nil
		}

		version := PackageVersion{
			Version:   pkg.Version,
			CreatedAt: timestamp.UTC(),
		}

		json := jsoniter.ConfigCompatibleWithStandardLibrary

		data, err := json.Marshal(version)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrMarshalVersionData, err)
		}

		if err = pkgBucket.Put(pkgVersion, data); err != nil {
			return fmt.Errorf("%w: %w", ErrStoreVersion, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTrackVersion, err)
	}

	return nil
}

// GetLastUpdate returns the last update time for a package.
func (m *Manager) GetLastUpdate(pkg *config.Package, timestamp time.Time) time.Time {
	var lastUpdate time.Time

	err := m.db.DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(xunsafe.StringToBytes(database.BucketPackageHistory))
		if bucket == nil {
			return fmt.Errorf("%w %q: %w", ErrGetHistory, pkg.Name, ErrBucketNotFound)
		}

		pkgBucket := bucket.Bucket(xunsafe.StringToBytes(pkg.Name))
		if pkgBucket == nil {
			return nil
		}

		cursor := pkgBucket.Cursor()

		for key, value := cursor.Last(); key != nil; key, value = cursor.Prev() {
			var version PackageVersion

			if err := json.Unmarshal(value, &version); err != nil {
				return fmt.Errorf("%w: %w: %w", ErrGetHistory, ErrMarshalVersionData, err)
			}

			if version.CreatedAt.After(lastUpdate) {
				lastUpdate = version.CreatedAt
			}
		}

		return nil
	})
	if err != nil {
		return timestamp.UTC()
	}

	if lastUpdate.IsZero() {
		return timestamp.UTC()
	}

	return lastUpdate
}
