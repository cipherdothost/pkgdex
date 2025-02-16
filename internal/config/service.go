// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config

import (
	"net/mail"
	"net/url"
	"strings"

	"git.sr.ht/~jamesponddotco/credential-go"
	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/key"
	"go.cipher.host/pkgdex/internal/meta"
)

const (
	// ErrServiceMissingKey is the error returned when the service key is missing.
	ErrServiceMissingKey cmdkit.Error = "missing service key"

	// ErrServiceInvalidKey is the error returned when the service key is invalid.
	ErrServiceInvalidKey cmdkit.Error = "invalid service key; must be at least 64 characters long and start with 'pkgdex_'"

	// ErrServiceInvalidContact is the error returned when the contact information is invalid.
	ErrServiceInvalidContact cmdkit.Error = "invalid contact information"

	// ErrServiceInvalidHomepage is the error returned when the homepage URL is invalid.
	ErrServiceInvalidHomepage cmdkit.Error = "invalid homepage URL"

	// ErrServiceInvalidBaseURL is the error returned when the base URL is invalid.
	ErrServiceInvalidBaseURL cmdkit.Error = "invalid base URL"

	// ErrServiceInvalidPrivacyPolicy is the error returned when the privacy policy URL is invalid.
	ErrServiceInvalidPrivacyPolicy cmdkit.Error = "invalid privacy policy URL"
)

const defaultServiceKeyLength int = 64

// Default service configuration values.
const (
	defaultServiceName        = meta.Name
	defaultServiceDescription = meta.LongDescription
	defaultServiceContact     = meta.Contact
	defaultServiceHomepage    = meta.Homepage
	defaultServiceBaseURL     = meta.BaseURL
	defaultPackagesPerPage    = 10
)

// SEO represents the service's SEO configuration.
type SEO struct {
	// Title is the title to use for the service's homepage. It will be suffixed
	// with the name of the service.
	Title string `json:"title"`

	// Description is the meta description to use for the service's homepage.
	Description string `json:"description"`

	// Image is the URL to an image which should represent the service.
	Image string `json:"image"`

	// ImageAlt is the description of the image.
	ImageAlt string `json:"imageAlt"`

	// Locale is the locale to use for the service.
	Locale string `json:"locale"`

	// Publisher is the Facebook page URL for the service.
	Publisher string `json:"publisher"`

	// Twitter is the Twitter handle for the service.
	Twitter string `json:"twitter"`
}

// Service represents the service configuration.
type Service struct {
	// SEO is the service's SEO configuration.
	SEO SEO `json:"seo"`

	// Name is the name of the service.
	Name string `json:"name"`

	// Description is a small description of the service to be displayed under
	// the name in the homepage.
	Description string `json:"description"`

	// Contact is the contact information for the service.
	Contact string `json:"contact"`

	// Homepage is the link to the service's homepage.
	Homepage string `json:"homepage"`

	// BaseURL is the base URL for package imports (e.g., "go.cipher.host/").
	BaseURL string `json:"baseURL"`

	// PrivacyPolicy is the link to the service's privacy policy.
	PrivacyPolicy string `json:"privacyPolicy"`

	// Key is an API key the service can use to protect sensitive routes.
	Key string

	// PackagesPerPage is the number of packages to display per page.
	PackagesPerPage int `json:"packagesPerPage"`

	// GeneratedKey tells the service if Key has been generated automatically.
	GeneratedKey bool

	// Archive tells the service whether to use the Wayback Machine to archive
	// pages.
	Archive bool `json:"archive"`
}

// validate validates the service configuration.
func (s *Service) validate() error {
	if s.Key == "" {
		return ErrServiceMissingKey
	}

	if len(s.Key) < defaultServiceKeyLength {
		return ErrServiceInvalidKey
	}

	if !strings.HasPrefix(s.Key, meta.Name+"_") {
		return ErrServiceInvalidKey
	}

	if _, err := mail.ParseAddress(s.Contact); err != nil {
		return ErrServiceInvalidContact
	}

	if _, err := url.Parse(s.Homepage); err != nil {
		return ErrServiceInvalidHomepage
	}

	if _, err := url.Parse(s.BaseURL); err != nil {
		return ErrServiceInvalidBaseURL
	}

	if s.PrivacyPolicy != "" {
		if _, err := url.Parse(s.PrivacyPolicy); err != nil {
			return ErrServiceInvalidPrivacyPolicy
		}
	}

	return nil
}

// setSecrets populates the configuration with the necessary secrets.
func (s *Service) setSecrets(store credential.Store) {
	secret, err := store.Get("key")
	if err != nil {
		secret = key.Generate()

		s.GeneratedKey = true
	}

	s.Key = strings.TrimSpace(secret)
}

// setDefaults sets the default values for the service configuration.
func (s *Service) setDefaults() {
	if s.Name == "" {
		s.Name = defaultServiceName
	}

	if s.Description == "" {
		s.Description = defaultServiceDescription
	}

	if s.Contact == "" {
		s.Contact = defaultServiceContact
	}

	if s.Homepage == "" {
		s.Homepage = defaultServiceHomepage
	}

	if s.BaseURL == "" {
		s.BaseURL = defaultServiceBaseURL
	}

	if s.PackagesPerPage == 0 {
		s.PackagesPerPage = defaultPackagesPerPage
	}
}
