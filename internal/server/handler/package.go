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
	"time"

	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/metric"
	"go.cipher.host/pkgdex/internal/repository"
	"go.cipher.host/pkgdex/internal/sanitize"
)

// packageData represents the data needed to render the package template.
type packageData struct {
	Service       *config.Service
	Usage         template.HTML
	Name          string
	ImportPath    string
	GoImport      string
	GoSource      string
	Description   string
	Version       string
	License       string
	Repository    string
	Documentation string
	Hidden        bool
}

// PackageHandler is the HTTP handler for the package endpoint.
type PackageHandler struct {
	logger   *slog.Logger
	metrics  *metric.Store
	tmpl     *template.Template
	service  *config.Service
	packages []*config.Package
}

// NewPackageHandler returns a new PackageHandler instance.
func NewPackageHandler(logger *slog.Logger, metrics *metric.Store, packages []*config.Package, service *config.Service, tmpl *template.Template) *PackageHandler {
	return &PackageHandler{
		logger:   logger,
		metrics:  metrics,
		service:  service,
		packages: packages,
		tmpl:     tmpl,
	}
}

// ServeHTTP handles HTTP requests for the packages endpoint.
func (h *PackageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pkgPath := strings.TrimPrefix(r.URL.Path, "/")
	pkgPath = strings.TrimSuffix(pkgPath, "/")

	var pkg *config.Package

	for _, p := range h.packages {
		if strings.HasPrefix(pkgPath, p.Name) {
			pkg = p

			break
		}
	}

	if pkg == nil {
		serve404(w, h.service, h.tmpl, pkgPath, h.logger)

		return
	}

	repo, err := repository.Parse(pkg.Repository, pkg.Branch)
	if err != nil {
		h.logger.Error("failed to parse repository",
			slog.String("repository", pkg.Repository),
			slog.String("branch", pkg.Branch),
			slog.Any("error", err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	// We want an approximate number of installs, so we only want to track installs from the
	// "go get" command. This can be spoofed and may not be 100% accurate, but it's the best
	// I can think of at the moment.
	if isGoGet(r) {
		if err = h.metrics.Track(pkg.Name, time.Now()); err != nil {
			// We log the error but continue serving the request because it makes the most sense.
			h.logger.Error("failed to track package installation",
				slog.String("package", pkg.Name),
				slog.Any("error", err),
			)
		}
	}

	importPath := path.Join(strings.TrimSuffix(h.service.BaseURL, "/"), pkgPath)

	var pkgUsage template.HTML

	if pkg.Usage != "" {
		pkgUsage, err = sanitize.Usage(pkg.Usage)
		if err != nil {
			h.logger.Error("failed to sanitize package usage",
				slog.String("package", pkg.Name),
				slog.Any("error", err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}
	}

	data := packageData{
		Service:       h.service,
		Name:          pkg.Name,
		ImportPath:    importPath,
		GoImport:      repo.GoImportMeta(importPath),
		GoSource:      repo.GoSourceMeta(importPath),
		Description:   pkg.Description,
		Version:       pkg.Version,
		License:       pkg.License,
		Repository:    pkg.Repository,
		Usage:         pkgUsage,
		Documentation: "https://pkg.go.dev/" + importPath,
		Hidden:        pkg.Hidden,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	setCacheHeaders(w, r)

	if err = h.tmpl.ExecuteTemplate(w, "package.html", data); err != nil {
		h.logger.Error("failed to execute template",
			slog.String("template", "package.html"),
			slog.Any("error", err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
