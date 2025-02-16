// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package handler

import (
	"log/slog"
	"net/http"

	"go.cipher.host/pkgdex/internal/sitemap"
)

// SitemapHandler is the HTTP handler for the sitemap.xml endpoint.
type SitemapHandler struct {
	logger *slog.Logger
	cache  *sitemap.Cache
}

// NewSitemapHandler returns a new SitemapHandler instance.
func NewSitemapHandler(logger *slog.Logger, cache *sitemap.Cache) *SitemapHandler {
	return &SitemapHandler{
		logger: logger,
		cache:  cache,
	}
}

// ServeHTTP handles HTTP requests for sitemap.xml endpoint.
func (h *SitemapHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	if _, err := w.Write(h.cache.Get()); err != nil {
		h.logger.Error("failed to write sitemap response",
			slog.Any("error", err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
