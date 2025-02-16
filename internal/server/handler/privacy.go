// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package handler

import (
	"html/template"
	"log/slog"
	"net/http"

	"go.cipher.host/pkgdex/internal/config"
)

// privacyData represents the data needed to render the privacy template.
type privacyData struct {
	Service *config.Service
}

// PrivacyHandler is the HTTP handler for the privacy endpoint.
type PrivacyHandler struct {
	logger  *slog.Logger
	service *config.Service
	tmpl    *template.Template
}

// NewPrivacyHandler returns a new PrivacyHandler instance.
func NewPrivacyHandler(logger *slog.Logger, service *config.Service, tmpl *template.Template) *PrivacyHandler {
	return &PrivacyHandler{
		logger:  logger,
		service: service,
		tmpl:    tmpl,
	}
}

// ServeHTTP handles HTTP requests for the privacy endpoint.
func (h *PrivacyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := privacyData{
		Service: h.service,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	setCacheHeaders(w, r)

	if err := h.tmpl.ExecuteTemplate(w, "privacy.html", data); err != nil {
		h.logger.Error("failed to execute template",
			slog.String("template", "privacy.html"),
			slog.Any("error", err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
