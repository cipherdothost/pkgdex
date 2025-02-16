// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

// Package key contains functions for generating API keys.
package key

import (
	"fmt"
	"hash/crc32"
	"strings"

	"git.sr.ht/~jamesponddotco/acopw-go"
	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
	"go.cipher.host/pkgdex/internal/meta"
)

const (
	separator       string = "_"
	randomKeyLength int    = 49
	keyLength       int    = 64
)

// Generate generates a new API key using the given prefix.
func Generate() string {
	var (
		prefix = meta.Name
		random = &acopw.Random{
			Length:     randomKeyLength,
			UseLower:   true,
			UseUpper:   true,
			UseNumbers: true,
			ExcludedCharset: []string{
				"i",
				"l",
				"o",
				"I",
				"L",
				"O",
			},
		}
		secret = random.Generate()
		key    = prefix + separator + secret
	)

	hasher := crc32.NewIEEE()
	hasher.Write(xunsafe.StringToBytes(key))

	checksum := hasher.Sum32()

	key = fmt.Sprintf("%s%x", key, checksum)

	if len(key) < keyLength {
		key += strings.Repeat("0", keyLength-len(key))
	}

	return key[:keyLength]
}
