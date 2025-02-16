// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package static_test

import (
	"strings"
	"testing"

	"go.cipher.host/pkgdex/internal/testhelper"
	"go.cipher.host/pkgdex/static"
)

func TestVersionHash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		givePath string
		wantLen  int
		wantErr  bool
	}{
		{
			name:     "valid scss file",
			givePath: "assets/scss/main.scss",
			wantLen:  16,
			wantErr:  false,
		},
		{
			name:     "highlight scss file",
			givePath: "assets/scss/highlight.scss",
			wantLen:  16,
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			givePath: "does-not-exist.css",
			wantLen:  0,
			wantErr:  true,
		},
		{
			name:     "empty path",
			givePath: "",
			wantLen:  0,
			wantErr:  true,
		},
		{
			name:     "path with special characters",
			givePath: "assets/css/#main%.css",
			wantLen:  0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := static.VersionHash(tt.givePath)
			if tt.wantErr {
				if err == nil {
					t.Error("VersionHash() error = nil, wantErr = true")
				}

				return
			}

			if err != nil {
				t.Errorf("VersionHash() error = %v, wantErr = false", err)

				return
			}

			if len(got) != tt.wantLen {
				t.Errorf("VersionHash() hash length = %d, want %d", len(got), tt.wantLen)
			}

			if !testhelper.IsHex(got) {
				t.Errorf("VersionHash() result is not valid hex: %v", got)
			}
		})
	}
}

func TestVersionHash_Consistency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		givePath string
	}{
		{
			name:     "main scss consistency",
			givePath: "assets/scss/main.scss",
		},
		{
			name:     "highlight scss consistency",
			givePath: "assets/scss/highlight.scss",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hash1, err := static.VersionHash(tt.givePath)
			if err != nil {
				t.Fatalf("First VersionHash() call failed: %v", err)
			}

			hash2, err := static.VersionHash(tt.givePath)
			if err != nil {
				t.Fatalf("Second VersionHash() call failed: %v", err)
			}

			if hash1 != hash2 {
				t.Errorf("VersionHash() not consistent: first = %v, second = %v", hash1, hash2)
			}
		})
	}
}

func TestVersionedPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveFilepath string
		wantPrefix   string
		wantExt      string
		wantErr      bool
	}{
		{
			name:         "valid scss file",
			giveFilepath: "assets/scss/main.scss",
			wantPrefix:   "assets/scss/main.",
			wantExt:      ".scss",
			wantErr:      false,
		},
		{
			name:         "highlight scss file",
			giveFilepath: "assets/scss/highlight.scss",
			wantPrefix:   "assets/scss/highlight.",
			wantExt:      ".scss",
			wantErr:      false,
		},
		{
			name:         "non-existent file",
			giveFilepath: "does-not-exist.css",
			wantErr:      true,
		},
		{
			name:         "empty path",
			giveFilepath: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := static.VersionedPath(tt.giveFilepath)
			if tt.wantErr {
				if err == nil {
					t.Error("VersionedPath() error = nil, wantErr = true")
				}

				return
			}

			if err != nil {
				t.Errorf("VersionedPath() error = %v, wantErr = false", err)

				return
			}

			if !strings.HasPrefix(got, tt.wantPrefix) {
				t.Errorf("VersionedPath() got = %v, want prefix %v", got, tt.wantPrefix)
			}

			if !strings.HasSuffix(got, tt.wantExt) {
				t.Errorf("VersionedPath() got = %v, want suffix %v", got, tt.wantExt)
			}

			hashPart := strings.TrimPrefix(got, tt.wantPrefix)
			hashPart = strings.TrimSuffix(hashPart, tt.wantExt)

			if len(hashPart) != 16 {
				t.Errorf("VersionedPath() hash length = %d, want 16", len(hashPart))
			}

			if !testhelper.IsHex(hashPart) {
				t.Errorf("VersionedPath() hash is not valid hex: %v", hashPart)
			}
		})
	}
}
