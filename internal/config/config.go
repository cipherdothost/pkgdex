// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

// Package config implements the configuration logic for the service.
package config

import (
	"fmt"
	"os"

	"git.sr.ht/~jamesponddotco/credential-go"
	jsoniter "github.com/json-iterator/go"
	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/meta"
)

const (
	// ErrInvalidConfigFile is the error returned when the configuration file is
	// invalid for whatever reason.
	ErrInvalidConfigFile cmdkit.Error = "invalid configuration file"

	// ErrConfigMissingService is the error returned when the service
	// configuration is missing.
	ErrConfigMissingService cmdkit.Error = "missing service configuration"

	// ErrConfigMissingServer is the error returned when the server
	// configuration is missing.
	ErrConfigMissingServer cmdkit.Error = "missing server configuration"

	// ErrConfigMissingDatabase is the error returned when the database
	// configuration is missing.
	ErrConfigMissingDatabase cmdkit.Error = "missing database configuration"

	// ErrConfigMissingPackages is the error returned when the packages
	// configuration is missing.
	ErrConfigMissingPackages cmdkit.Error = "missing package configuration"
)

// DefaultConfigLocation is the default location for the configuration file.
const DefaultConfigLocation = "/etc/" + meta.Name + "/config.json"

// Config represents the configuration for the service.
type Config struct {
	// Service is the service configuration.
	Service *Service `json:"service"`

	// Server is the server configuration.
	Server *Server `json:"server"`

	// Database is the database configuration.
	Database *Database `json:"database"`

	// Packages is the package configuration.
	Packages []*Package `json:"packages"`
}

// Load loads the configuration from the given path.
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}
	defer file.Close()

	cfg := &Config{
		Service:  &Service{},
		Server:   &Server{},
		Database: &Database{},
		Packages: make([]*Package, 0),
	}

	cfg.setDefaults()

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}

	if err = cfg.setSecrets(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}

	if err = cfg.validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}

	return cfg, nil
}

// validate goes through the configuration and ensures it is valid.
func (c *Config) validate() error {
	if c.Service == nil {
		return ErrConfigMissingService
	}

	if err := c.Service.validate(); err != nil {
		return fmt.Errorf("%w", err)
	}

	if c.Server == nil {
		return ErrConfigMissingServer
	}

	if err := c.Server.validate(); err != nil {
		return fmt.Errorf("%w", err)
	}

	if c.Database == nil {
		return ErrConfigMissingDatabase
	}

	if err := c.Database.validate(); err != nil {
		return fmt.Errorf("%w", err)
	}

	if c.Packages == nil {
		return ErrConfigMissingPackages
	}

	for _, pkg := range c.Packages {
		if err := pkg.validate(); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

// setSecrets populates the configuration with the given secrets.
func (c *Config) setSecrets() error {
	store, err := credential.Open(meta.Name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = c.Server.setSecrets(store); err != nil {
		return fmt.Errorf("%w", err)
	}

	c.Service.setSecrets(store)

	return nil
}

// setDefaults populates the configuration with the default values.
func (c *Config) setDefaults() {
	c.Service.setDefaults()
	c.Server.setDefaults()
	c.Database.setDefaults()
}
