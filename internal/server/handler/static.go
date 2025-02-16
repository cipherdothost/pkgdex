// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package handler

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"slices"
	"strings"
	"sync"

	"go.cipher.host/pkgdex/internal/config"
)

// defaultCacheCapacity is the default capacity used to initialize the cache
// map. It's the value of the known files in the static directory in bytes, plus
// a few extra for safety.
const defaultCacheCapacity = 18000

// StaticHandler is the HTTP handler for the static files.
type StaticHandler struct {
	tmpl      *template.Template
	service   *config.Service
	logger    *slog.Logger
	files     embed.FS
	cache     map[string][]byte
	cacheLock sync.RWMutex
}

// NewStaticHandler returns a new StaticHandler instance.
func NewStaticHandler(logger *slog.Logger, tmpl *template.Template, service *config.Service, files embed.FS) *StaticHandler {
	return &StaticHandler{
		tmpl:    tmpl,
		service: service,
		logger:  logger,
		files:   files,
		cache:   make(map[string][]byte, defaultCacheCapacity),
	}
}

// ServeHTTP handles HTTP requests for static files.
func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		fileName  = path.Clean(r.URL.Path[1:])
		extension = path.Ext(fileName)
	)

	if extension == ".css" {
		var (
			dir, file = path.Split(fileName)
			parts     = strings.Split(file, ".")
		)

		if len(parts) == 3 {
			fileName = path.Join(dir, parts[0]+".css")
		}
	}

	allowedPaths := []string{
		"favicon.ico",
		"robots.txt",
		"site.webmanifest",
		"assets/css/main.css",
		"assets/css/highlight.css",
		"assets/img/android-chrome-192x192.png",
		"assets/img/android-chrome-512x512.png",
		"assets/img/apple-touch-icon.png",
		"assets/img/favicon-16x16.png",
		"assets/img/favicon-32x32.png",
		"assets/img/favicon.svg",
	}

	if !slices.Contains(allowedPaths, fileName) {
		h.logger.Error("static file not allowed",
			slog.String("file", fileName),
		)

		serve404(w, h.service, h.tmpl, r.URL.Path, h.logger)

		return
	}

	content, err := h.getContent(fileName)
	if err != nil {
		h.logger.Error("failed to read static file",
			slog.String("file", fileName),
			slog.Any("error", err),
		)

		serve404(w, h.service, h.tmpl, r.URL.Path, h.logger)

		return
	}

	var contentType string

	switch extension {
	case ".txt":
		contentType = "text/plain; charset=utf-8"

		w.Header().Set("Cache-Control", "public, max-age=86400")
	case ".webmanifest":
		contentType = "application/manifest+json; charset=utf-8"

		w.Header().Set("Cache-Control", "public, max-age=2592000")
	case ".css":
		contentType = "text/css; charset=utf-8"

		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	case ".png":
		contentType = "image/png"

		w.Header().Set("Cache-Control", "public, max-age=2592000")
	case ".ico":
		contentType = "image/x-icon"

		w.Header().Set("Cache-Control", "public, max-age=2592000")
	case ".svg":
		contentType = "image/svg+xml; charset=utf-8"

		w.Header().Set("Cache-Control", "public, max-age=2592000")
	default:
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)

	if _, err = w.Write(content); err != nil {
		h.logger.Error("failed to write static file response",
			slog.String("file", fileName),
			slog.Any("error", err),
		)
	}
}

// getContent retrieves file content from cache or filesystem.
func (h *StaticHandler) getContent(fileName string) ([]byte, error) {
	h.cacheLock.RLock()
	if content, ok := h.cache[fileName]; ok {
		h.cacheLock.RUnlock()

		return content, nil
	}
	h.cacheLock.RUnlock()

	content, err := fs.ReadFile(h.files, fileName)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if fileName == "robots.txt" {
		content = bytes.ReplaceAll(content, []byte("{{SERVICE_CONTACT}}"), []byte(h.service.Contact))
		content = bytes.ReplaceAll(content, []byte("{{SERVICE_BASEURL}}"), []byte(h.service.BaseURL))
	}

	h.cacheLock.Lock()
	h.cache[fileName] = content
	h.cacheLock.Unlock()

	return content, nil
}
