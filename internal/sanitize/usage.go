// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package sanitize

import (
	"html/template"

	"github.com/microcosm-cc/bluemonday"
	"go.cipher.host/cmdkit"
)

// ErrEmptyUsageString is returned when the sanitized usage string is empty.
const ErrEmptyUsageString cmdkit.Error = "usage example is empty after sanitization; ensure the input contains valid content"

// Usage tries to safely sanitizes HTML content for the config.Package.Usage
// field to prevent XSS attacks.
func Usage(input string) (template.HTML, error) {
	policy := bluemonday.NewPolicy()

	// Allow necessary elements
	policy.AllowElements("pre", "code", "span", "a")

	// Allow attributes on pre
	policy.AllowAttrs("class", "tabindex").OnElements("pre")

	// Allow attributes on code
	policy.AllowAttrs("class").OnElements("code")

	// Allow attributes on span
	policy.AllowAttrs("class", "id").OnElements("span")

	// Allow attributes on anchors
	policy.AllowAttrs("style", "href", "id").OnElements("a")

	// Allow specific styles
	policy.AllowStyles("outline", "text-decoration", "color").Globally()

	// Sanitize the input
	sanitized := policy.Sanitize(input)
	if sanitized == "" {
		return "", ErrEmptyUsageString
	}

	return template.HTML(sanitized), nil //nolint:gosec // we just sanitized the input
}
