// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package testhelper

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go.cipher.host/pkgdex/internal/meta"
)

// SetupSecrets writes the provided secrets to files in the given directory.
func SetupSecrets(t *testing.T, dir string, secrets map[string]string) error {
	t.Helper()

	for k, v := range secrets {
		if err := os.WriteFile(filepath.Join(dir, meta.Name+"-"+k), []byte(v), 0o600); err != nil {
			return fmt.Errorf("failed to write secret %q: %w", k, err)
		}
	}

	return nil
}
