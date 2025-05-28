// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package handler

import (
	"bytes"
	"log/slog"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/metric"
	"go.cipher.host/pkgdex/internal/server/model"
	"go.cipher.host/x/xnet/xhttp"
)

// DownloadsHandler is the HTTP handler for the /meta/downloads endpoint.
type DownloadsHandler struct {
	logger   *slog.Logger
	metrics  *metric.Store
	packages []*config.Package
}

// NewDownloadsHandler returns a new DownloadsHandler instance.
func NewDownloadsHandler(logger *slog.Logger, metrics *metric.Store, packages []*config.Package) *DownloadsHandler {
	return &DownloadsHandler{
		logger:   logger,
		metrics:  metrics,
		packages: packages,
	}
}

// ServeHTTP handles HTTP requests for the /meta/downloads endpoint.
func (d *DownloadsHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	report := &model.Report{
		Packages: make([]model.Package, 0, len(d.packages)),
	}

	for _, pkg := range d.packages {
		count, err := d.metrics.Count(pkg.Name)
		if err != nil {
			d.logger.Error("failed to get download count",
				slog.String("package", pkg.Name),
				slog.Any("error", err),
			)

			continue
		}

		report.Packages = append(report.Packages, model.Package{
			Name:      pkg.Name,
			Downloads: int(count),
		})
	}

	var (
		jsonBuf bytes.Buffer
		json    = jsoniter.ConfigCompatibleWithStandardLibrary
		encoder = json.NewEncoder(&jsonBuf)
	)

	w.Header().Set("Content-Type", "application/json")

	if err := encoder.Encode(report); err != nil {
		d.logger.Error("failed to encode download report response",
			slog.Any("error", err),
		)

		response := xhttp.DetailError{
			Detail: "Failed to encode download report response.",
			Status: http.StatusInternalServerError,
		}

		if err = response.WriteJSON(w); err != nil {
			d.logger.Error("failed to write error response",
				slog.Any("error", err),
			)
		}

		return
	}

	if _, err := w.Write(jsonBuf.Bytes()); err != nil {
		d.logger.Error("failed to write download report response",
			slog.Any("error", err),
		)
	}
}
