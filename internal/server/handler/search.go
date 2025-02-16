// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package handler

import (
	"html/template"
	"log/slog"
	"net/http"

	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/pagination"
	"go.cipher.host/pkgdex/internal/search"
)

// searchData represents the data needed to render the search template.
type searchData struct {
	Service    *config.Service
	Query      string
	Results    []search.Package
	Pagination pagination.Info
	NoResults  bool
}

// SearchHandler handles search requests.
type SearchHandler struct {
	logger  *slog.Logger
	manager *search.Manager
	service *config.Service
	tmpl    *template.Template
}

// NewSearchHandler creates a new search handler.
func NewSearchHandler(logger *slog.Logger, manager *search.Manager, service *config.Service, tmpl *template.Template) *SearchHandler {
	return &SearchHandler{
		logger:  logger,
		manager: manager,
		service: service,
		tmpl:    tmpl,
	}
}

// ServeHTTP handles HTTP requests for the search endpoint.
func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	results, err := h.manager.Search(query)
	if err != nil {
		h.logger.Error("search failed",
			slog.String("query", query),
			slog.Any("error", err),
		)

		http.Error(w, "Search failed", http.StatusInternalServerError)

		return
	}

	var (
		currentPage    = pagination.GetPage(r)
		paginationInfo = pagination.New(len(results), currentPage, h.service.PackagesPerPage)
		start          = paginationInfo.Offset()
		end            = start + paginationInfo.Limit()
	)

	if end > len(results) {
		end = len(results)
	}

	pageResults := results[start:end]

	data := searchData{
		Service:    h.service,
		Query:      query,
		Results:    pageResults,
		NoResults:  len(results) == 0,
		Pagination: paginationInfo,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	setCacheHeaders(w, r)

	if err = h.tmpl.ExecuteTemplate(w, "search.html", data); err != nil {
		h.logger.Error("failed to execute template",
			slog.String("template", "search.html"),
			slog.Any("error", err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
