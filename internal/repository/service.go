// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package repository

// Service represents a version control system service or provider.
type Service string

// Supported VCS services.
const (
	ServiceGitHub    Service = "github.com"
	ServiceSourcehut Service = "git.sr.ht"
	ServiceGitLab    Service = "gitlab.com"
)
