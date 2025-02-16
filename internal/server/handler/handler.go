// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

// Package handler contains the HTTP handlers for the service.
package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"go.cipher.host/pkgdex/internal/config"
)

// setCacheHeaders sets the HTML cache headers for the response.
func setCacheHeaders(w http.ResponseWriter, r *http.Request) {
	if isGoGet(r) {
		w.Header().Set("Cache-Control", "public, max-age=86400")

		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
}

// isGoGet returns true if the request is coming from the "go get" command.
func isGoGet(r *http.Request) bool {
	var (
		userAgent = r.Header.Get("User-Agent")
		goGet     = r.URL.Query().Get("go-get") == "1"
	)

	return strings.HasPrefix(userAgent, "Go-http-client") && goGet
}

// serve404 serves a 404 page.
func serve404(w http.ResponseWriter, service *config.Service, tmpl *template.Template, path string, logger *slog.Logger) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")

	w.WriteHeader(http.StatusNotFound)

	data := struct {
		Service *config.Service
		Path    string
	}{
		Service: service,
		Path:    path,
	}

	if err := tmpl.ExecuteTemplate(w, "404.html", data); err != nil {
		logger.Error("failed to execute template",
			slog.String("template", "404.html"),
			slog.Any("error", err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
