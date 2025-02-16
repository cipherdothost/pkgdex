// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package database

import (
	"fmt"

	"go.cipher.host/cmdkit"
	bolt "go.etcd.io/bbolt"
)

const (
	// ErrHealthCheck is returned when a health check fails.
	ErrHealthCheck cmdkit.Error = "database health check failed; the database may be corrupted or inaccessible"

	// ErrNilBucket is returned when a database bucket is nil.
	ErrNilBucket cmdkit.Error = "database bucket not found; database may be corrupted or incorrectly initialize"
)

// HealthCheck performs a health check on the database.
func (s *Store) HealthCheck() error {
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			if b == nil {
				return fmt.Errorf("%w: %q", ErrNilBucket, name)
			}

			return nil
		})
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrHealthCheck, err)
	}

	return nil
}
