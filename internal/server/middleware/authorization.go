// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package middleware

import (
	"log/slog"
	"net/http"

	"go.cipher.host/x/xcrypto/xsubtle"
	"go.cipher.host/x/xnet/xhttp"
)

// Authorization ensures that the request has a valid API key.
func Authorization(apiKey string, logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !xsubtle.ConstantTimeCompareString(r.Header.Get("Authorization"), "Bearer "+apiKey) {
			logger.Warn("unauthorized API access attempt",
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("path", r.URL.Path),
				slog.String("key", r.Header.Get("Authorization")),
			)

			response := xhttp.DetailError{
				Status: http.StatusUnauthorized,
				Detail: "Access denied. Please provide a valid API key in the Authorization header.",
			}

			if err := response.WriteJSON(w); err != nil {
				logger.Error("failed to write response",
					slog.Any("error", err),
				)
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}
