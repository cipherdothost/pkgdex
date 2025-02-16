// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config

import (
	"net/url"
	"regexp"

	"go.cipher.host/cmdkit"
)

const (
	validPackageNameRegex    = `^[a-zA-Z0-9][a-zA-Z0-9-_/.]*[a-zA-Z0-9]$`
	validPackageVersionRegex = `^\d+\.\d+\.\d+$`
)

const (
	// ErrPackageMissingName is returned when the package name is missing from the
	// configuration.
	ErrPackageMissingName cmdkit.Error = "package name is missing"

	// ErrPackageMissingDescription is returned when the package description is
	// missing from the configuration.
	ErrPackageMissingDescription cmdkit.Error = "package description is missing"

	// ErrPackageMissingVersion is returned when the package version is missing
	// from the configuration.
	ErrPackageMissingVersion cmdkit.Error = "package version is missing"

	// ErrPackageMissingBranch is returned when the package branch is missing
	// from the configuration.
	ErrPackageMissingBranch cmdkit.Error = "package branch is missing"

	// ErrPackageMissingRepository is returned when the package repository is missing
	// from the configuration.
	ErrPackageMissingRepository cmdkit.Error = "package repository is missing"

	// ErrPackageMissingLicense is returned when the package license is missing
	// from the configuration.
	ErrPackageMissingLicense cmdkit.Error = "package license is missing"

	// ErrPackageInvalidName is returned when the package name is invalid.
	ErrPackageInvalidName cmdkit.Error = "package name is invalid"

	// ErrPackageInvalidVersion is returned when the package version is invalid.
	ErrPackageInvalidVersion cmdkit.Error = "package version is invalid; must be in the format x.y.z"

	// ErrPackageInvalidRepository is returned when the package repository URL is invalid.
	ErrPackageInvalidRepository cmdkit.Error = "package repository URL is invalid"
)

// Package represents a package configuration.
type Package struct {
	// Name is the name of the package.
	Name string `json:"name"`

	// Description is the description of the package.
	Description string `json:"description"`

	// Version is the version of the package.
	Version string `json:"version"`

	// License is the license for the package.
	License string `json:"license"`

	// Repository is the repository for the package.
	Repository string `json:"repository"`

	// Branch is the branch for the package.
	Branch string `json:"branch"`

	// Usage is a piece of example code showing how to use the package.
	Usage string `json:"usage"`

	// Hidden tells the service that this package should not be shown to users
	// in the index, search results, sitemap, etc.
	Hidden bool `json:"hidden"`
}

// validate validates the package configuration.
func (p *Package) validate() error {
	if p.Name == "" {
		return ErrPackageMissingName
	}

	nameRegex := regexp.MustCompile(validPackageNameRegex)

	if !nameRegex.MatchString(p.Name) {
		return ErrPackageInvalidName
	}

	if p.Description == "" {
		return ErrPackageMissingDescription
	}

	if p.Version == "" {
		return ErrPackageMissingVersion
	}

	versionRegex := regexp.MustCompile(validPackageVersionRegex)

	if !versionRegex.MatchString(p.Version) {
		return ErrPackageInvalidVersion
	}

	if p.Branch == "" {
		return ErrPackageMissingBranch
	}

	if p.Repository == "" {
		return ErrPackageMissingRepository
	}

	if _, err := url.Parse(p.Repository); err != nil {
		return ErrPackageInvalidRepository
	}

	if p.License == "" {
		return ErrPackageMissingLicense
	}

	return nil
}
