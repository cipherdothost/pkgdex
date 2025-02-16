// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package handler

import (
	"log/slog"
	"net/http"

	"go.cipher.host/pkgdex/internal/rss"
)

// FeedHandler is the HTTP handler for the feed.xml endpoint.
type FeedHandler struct {
	logger *slog.Logger
	cache  *rss.Cache
}

// NewFeedHandler returns a new FeedHandler instance.
func NewFeedHandler(logger *slog.Logger, cache *rss.Cache) *FeedHandler {
	return &FeedHandler{
		logger: logger,
		cache:  cache,
	}
}

// ServeHTTP handles HTTP requests for feed.xml endpoint.
func (h *FeedHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	if _, err := w.Write(h.cache.Get()); err != nil {
		h.logger.Error("failed to write RSS feed response",
			slog.Any("error", err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
