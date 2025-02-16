// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package middleware

import (
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xstrings"
)

// CSP ensures that the Content-Security-Policy header is set.
func CSP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			directives = []string{
				"default-src 'self'",
				"script-src 'none'",
				"script-src-elem 'none'",
				"script-src-attr 'none'",
				"style-src 'self' 'unsafe-inline'",
				"style-src-elem 'self'",
				"style-src-attr 'self' 'unsafe-inline'",
				"img-src 'self' data:",
				"font-src 'none'",
				"connect-src 'none'",
				"media-src 'none'",
				"object-src 'none'",
				"child-src 'none'",
				"frame-src 'none'",
				"worker-src 'none'",
				"base-uri 'self'",
				"manifest-src 'self'",
				"form-action 'self'",
				"frame-ancestors 'none'",
				"block-all-mixed-content",
				"upgrade-insecure-requests",
			}
			csp = xstrings.JoinWithSeparator("; ", directives...)
		)

		w.Header().Set("Content-Security-Policy", csp)

		next.ServeHTTP(w, r)
	})
}
