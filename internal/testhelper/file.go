// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package testhelper

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// CopyConfigFiles copies the configuration files from the source to the
// destination directory, including subdirectories. It'll ignore any
// "config.json" files encountered during traversal.
func CopyConfigFiles(t *testing.T, src, dst string) error {
	t.Helper()

	if err := os.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("creating destination directory: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("reading source directory: %w", err)
	}

	for _, entry := range entries {
		var (
			srcPath = filepath.Join(src, entry.Name())
			dstPath = filepath.Join(dst, entry.Name())
		)

		if entry.Name() == "config.json" {
			continue
		}

		var info fs.FileInfo

		info, err = entry.Info()
		if err != nil {
			return fmt.Errorf("getting file info for %q: %w", entry.Name(), err)
		}

		if info.IsDir() {
			if err = CopyConfigFiles(t, srcPath, dstPath); err != nil {
				return fmt.Errorf("copying directory %q: %w", entry.Name(), err)
			}

			continue
		}

		if err = copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("copying file %q: %w", entry.Name(), err)
		}
	}

	return nil
}

// copyFile copies a single file from source to destination.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("creating destination file: %w", err)
	}

	defer func() {
		closeErr := dstFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copying file contents: %w", err)
	}

	return nil
}
