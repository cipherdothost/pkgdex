// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

// Package testhelper provides functions and utilities for testing.
package testhelper

// IsHex returns true if the string contains only hexadecimal characters.
func IsHex(str string) bool {
	for _, rune := range str {
		if !((rune >= '0' && rune <= '9') || (rune >= 'a' && rune <= 'f') || (rune >= 'A' && rune <= 'F')) {
			return false
		}
	}

	return true
}
