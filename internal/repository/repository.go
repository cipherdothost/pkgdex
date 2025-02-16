// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package repository

import (
	"fmt"
	"net/url"
	"strings"

	"go.cipher.host/cmdkit"
)

const (
	// ErrInvalidURL is returned when the repository URL is invalid.
	ErrInvalidURL cmdkit.Error = "invalid repository URL; please provide a valid HTTPS URL"

	// ErrMissingURL is returned when the repository URL is missing.
	ErrMissingURL cmdkit.Error = "repository URL is required; please specify the source repository URL"

	// ErrMissingBranch is returned when the repository branch is missing.
	ErrMissingBranch cmdkit.Error = "repository branch is required; please specify the source code branch"

	// ErrUnsupportedService is returned when the repository service is not supported.
	ErrUnsupportedService cmdkit.Error = "unsupported version control service; supported services include GitHub, GitLab, and Sourcehut"
)

// Repository represents a VCS repository.
type Repository struct {
	// Service is the VCS service or provider (e.g., ServiceGitHub, ServiceSourcehut).
	Service Service

	// Owner is the repository owner or organization.
	Owner string

	// Name is the repository name.
	Name string

	// Branch is the repository branch.
	Branch string
}

// Parse parses a repository URL and returns a Repository instance with the specified branch.
func Parse(repository, branch string) (*Repository, error) {
	if repository == "" {
		return nil, ErrMissingURL
	}

	if branch == "" {
		return nil, ErrMissingBranch
	}

	if !strings.Contains(repository, "://") {
		repository = "https://" + repository
	}

	u, err := url.Parse(repository)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	var (
		host    = strings.TrimPrefix(u.Host, "www.")
		service = Service(host)
	)

	if service != ServiceGitHub && service != ServiceGitLab && service != ServiceSourcehut {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedService, service)
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, repository)
	}

	return &Repository{
		Service: Service(host),
		Owner:   parts[0],
		Name:    parts[1],
		Branch:  branch,
	}, nil
}

// CloneURL returns the HTTPS clone URL for the repository.
func (r *Repository) CloneURL() string {
	return fmt.Sprintf("https://%s/%s/%s", r.Service, r.Owner, r.Name)
}

// SourceURL returns the URL to view source code in the repository.
func (r *Repository) SourceURL(file, line string) string {
	switch r.Service {
	case ServiceSourcehut:
		if line != "" {
			return fmt.Sprintf("https://%s/%s/%s/tree/%s/item/%s#L%s", r.Service, r.Owner, r.Name, r.Branch, file, line)
		}

		return fmt.Sprintf("https://%s/%s/%s/tree/%s/item/%s", r.Service, r.Owner, r.Name, r.Branch, file)
	case ServiceGitHub:
		if line != "" {
			return fmt.Sprintf("https://%s/%s/%s/blob/%s/%s#L%s", r.Service, r.Owner, r.Name, r.Branch, file, line)
		}

		return fmt.Sprintf("https://%s/%s/%s/blob/%s/%s", r.Service, r.Owner, r.Name, r.Branch, file)
	case ServiceGitLab:
		if line != "" {
			return fmt.Sprintf("https://%s/%s/%s/-/blob/%s/%s#L%s", r.Service, r.Owner, r.Name, r.Branch, file, line)
		}

		return fmt.Sprintf("https://%s/%s/%s/-/blob/%s/%s", r.Service, r.Owner, r.Name, r.Branch, file)
	default:
		return ""
	}
}

// GoImportMeta returns the go-import meta tag content.
func (r *Repository) GoImportMeta(importPath string) string {
	return fmt.Sprintf("%s git %s", importPath, r.CloneURL())
}

// GoSourceMeta returns the go-source meta tag content.
func (r *Repository) GoSourceMeta(importPath string) string {
	var (
		dirTemplate  string
		fileTemplate string
		uri          = r.CloneURL()
	)

	switch r.Service {
	case ServiceGitHub:
		dirTemplate = fmt.Sprintf("%s/tree/%s{/dir}", uri, r.Branch)
		fileTemplate = fmt.Sprintf("%s/blob/%s{/dir}/{file}#L{line}", uri, r.Branch)
	case ServiceGitLab:
		dirTemplate = fmt.Sprintf("%s/-/tree/%s{/dir}", uri, r.Branch)
		fileTemplate = fmt.Sprintf("%s/-/blob/%s{/dir}/{file}#L{line}", uri, r.Branch)
	case ServiceSourcehut:
		dirTemplate = fmt.Sprintf("%s/tree/%s/item{/dir}", uri, r.Branch)
		fileTemplate = fmt.Sprintf("%s/tree/%s/item{/dir}/{file}#L{line}", uri, r.Branch)
	}

	return fmt.Sprintf("%s %s %s %s", importPath, uri, dirTemplate, fileTemplate)
}
