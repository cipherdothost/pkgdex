// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/pagination"
)

// rootData represents the data needed to render the root template.
type rootData struct {
	Service    *config.Service
	Packages   []packageData
	Pagination pagination.Info
}

// RootHandler is the HTTP handler for the root endpoint.
type RootHandler struct {
	logger   *slog.Logger
	service  *config.Service
	tmpl     *template.Template
	packages []*config.Package
}

// NewRootHandler returns a new RootHandler instance.
func NewRootHandler(logger *slog.Logger, service *config.Service, packages []*config.Package, tmpl *template.Template) *RootHandler {
	return &RootHandler{
		logger:   logger,
		service:  service,
		packages: packages,
		tmpl:     tmpl,
	}
}

// ServeHTTP handles HTTP requests for the root endpoint.
func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		serve404(w, h.service, h.tmpl, r.URL.Path, h.logger)

		return
	}

	visiblePackages := make([]*config.Package, 0, len(h.packages))

	for _, pkg := range h.packages {
		if !pkg.Hidden {
			visiblePackages = append(visiblePackages, pkg)
		}
	}

	var (
		currentPage    = pagination.GetPage(r)
		paginationInfo = pagination.New(len(visiblePackages), currentPage, h.service.PackagesPerPage)
		start          = paginationInfo.Offset()
		end            = start + paginationInfo.Limit()
	)

	if end > len(visiblePackages) {
		end = len(visiblePackages)
	}

	var (
		pagePackages = visiblePackages[start:end]
		packages     = make([]packageData, 0, len(pagePackages))
	)

	for _, pkg := range pagePackages {
		importPath := path.Join(strings.TrimSuffix(h.service.BaseURL, "/"), pkg.Name)

		packages = append(packages, packageData{
			Name:          pkg.Name,
			Description:   pkg.Description,
			License:       pkg.License,
			Repository:    pkg.Repository,
			Documentation: "https://pkg.go.dev/" + importPath,
		})
	}

	data := rootData{
		Service:    h.service,
		Packages:   packages,
		Pagination: paginationInfo,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	setCacheHeaders(w, r)

	if err := h.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		h.logger.Error("failed to execute template",
			slog.String("template", "index.html"),
			slog.Any("error", err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
