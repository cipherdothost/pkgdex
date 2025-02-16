// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package model

// Report represents a download statistics report.
type Report struct {
	Packages []Package `json:"packages"`
}

// Package represents download statistics for a package.
type Package struct {
	Name      string `json:"name"`
	Downloads int    `json:"downloads"`
}
