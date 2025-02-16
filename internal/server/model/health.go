// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package model

// Result represents the result of a health check.
type Result struct {
	Name         string `json:"name"`
	Label        string `json:"label"`
	Status       string `json:"status"`
	ShortSummary string `json:"shortSummary"`
}

// Health represents the health status of the service and its dependencies.
type Health struct {
	CheckResults []Result `json:"checkResults"`
	FinishedAt   int64    `json:"finishedAt"`
}
