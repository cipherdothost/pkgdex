// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package middleware

import (
	"log/slog"
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xsubtle"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

// Authorization ensures that the request has a valid API key.
func Authorization(apiKey string, logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !xsubtle.ConstantTimeStringEqual(r.Header.Get("Authorization"), "Bearer "+apiKey) {
			logger.Warn("unauthorized API access attempt",
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("path", r.URL.Path),
				slog.String("key", r.Header.Get("Authorization")),
			)

			response := xhttp.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: "Access denied. Please provide a valid API key in the Authorization header.",
			}

			response.Write(r.Context(), logger, w)

			return
		}

		next.ServeHTTP(w, r)
	})
}
