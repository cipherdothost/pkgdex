// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package handler

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.cipher.host/pkgdex/internal/meta"
	"go.cipher.host/pkgdex/internal/server/model"
	"go.cipher.host/x/xnet/xhttp"
)

const (
	// ok is the status of a service that is okay.
	ok string = "ok"
)

// HealthHandler is the HTTP handler for the /health endpoint.
type HealthHandler struct {
	logger *slog.Logger
}

// NewHealthHandler returns a new HealthHandler instance.
func NewHealthHandler(logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// ServeHTTP handles HTTP requests for the /health endpoint.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	health := model.Health{
		FinishedAt: time.Now().Unix(),
		CheckResults: []model.Result{
			{
				Name:         "VersionCheck",
				Label:        "Version check",
				Status:       ok,
				ShortSummary: "Running version " + meta.Version + " of " + meta.Name + ".",
			},
		},
	}

	var (
		jsonBuf bytes.Buffer
		json    = jsoniter.ConfigCompatibleWithStandardLibrary
		encoder = json.NewEncoder(&jsonBuf)
	)

	w.Header().Set("Content-Type", "application/json")

	if err := encoder.Encode(health); err != nil {
		h.logger.Error("failed to encode health response",
			slog.Any("error", err),
		)

		response := xhttp.DetailError{
			Detail: "Failed to encode health response.",
			Status: http.StatusInternalServerError,
		}

		response.WriteJSON(w)

		return
	}

	if _, err := w.Write(jsonBuf.Bytes()); err != nil {
		h.logger.Error("failed to write health response",
			slog.Any("error", err),
		)
	}
}
