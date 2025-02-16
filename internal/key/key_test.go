// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package key_test

import (
	"fmt"
	"hash/crc32"
	"strings"
	"testing"

	"go.cipher.host/pkgdex/internal/key"
	"go.cipher.host/pkgdex/internal/meta"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	got := key.Generate()

	if want := 64; len(got) != want {
		t.Errorf("Generate() returned key with length %d, want %d", len(got), want)
	}

	if !strings.HasPrefix(got, meta.Name+"_") {
		t.Errorf("Generate() returned key %q that doesn't start with %q_", got, meta.Name)
	}

	parts := strings.Split(got, "_")
	if len(parts) != 2 {
		t.Fatalf("Generate() returned malformed key %q", got)
	}

	randomPart := parts[1]
	if len(randomPart) < 57 {
		t.Fatalf("Generate() returned key with invalid random part %q", randomPart)
	}

	checksumStr := randomPart[len(randomPart)-8:]
	keyWithoutChecksum := got[:len(got)-8]

	hasher := crc32.NewIEEE()
	hasher.Write([]byte(keyWithoutChecksum))
	want := hasher.Sum32()

	if got = checksumStr; got != strings.ToLower(fmt.Sprintf("%08x", want)) {
		t.Errorf("Generate() has invalid checksum %q, want %08x", got, want)
	}
}

func TestGenerate_Uniqueness(t *testing.T) {
	t.Parallel()

	const iterations = 1000

	seen := make(map[string]struct{}, iterations)

	for range iterations {
		got := key.Generate()

		if _, exists := seen[got]; exists {
			t.Errorf("Generate() produced duplicate key: %q", got)
		}

		seen[got] = struct{}{}
	}
}
