// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package wayback

import "git.sr.ht/~jamesponddotco/xstd-go/xstrings"

// PackageURI returns the URI of the given package.
func PackageURI(baseURL, pkg string) string {
	return xstrings.JoinWithSeparator("/", baseURL, pkg)
}
