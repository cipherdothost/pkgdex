// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package wayback

import (
	"fmt"
	"net/url"
)

// PackageURI returns the URI of the given package.
func PackageURI(baseURL, pkg string) (string, error) {
	uri, err := url.JoinPath(baseURL, pkg)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return uri, nil
}
