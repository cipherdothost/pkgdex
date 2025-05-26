// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config

import (
	"fmt"
	"time"

	"git.sr.ht/~jamesponddotco/credential-go"
	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/meta"
	"go.cipher.host/x/xtime"
)

const (
	// ErrServerMissingTLSCertificate is the error returned when the TLS certificate is missing.
	ErrServerMissingTLSCertificate cmdkit.Error = "missing TLS certificate"

	// ErrServerMissingTLSKey is the error returned when the TLS key is missing.
	ErrServerMissingTLSKey cmdkit.Error = "missing TLS key"

	// ErrServerMissingTLSConfig is the error returned when the TLS configuration is missing.
	ErrServerMissingTLSConfig cmdkit.Error = "missing TLS configuration"
)

// Default server configuration values.
const (
	defaultServerAddress      string = ":1997"
	defaultServerPID                 = "/var/run/" + meta.Name + "/" + meta.Name + ".pid"
	defaultServerReadTimeout         = xtime.Duration(5 * time.Second)
	defaultServerWriteTimeout        = xtime.Duration(10 * time.Second)
	defaultServerIdleTimeout         = xtime.Duration(30 * time.Second)
)

// TLS represents the TLS configuration.
type TLS struct {
	// Certificate is the path to the TLS certificate.
	Certificate []byte

	// Key is the path to the TLS key.
	Key []byte
}

// Server represents the basic server configuration.
type Server struct {
	// TLS is the TLS configuration.
	TLS *TLS

	// Address is the address to listen on.
	Address string `json:"address"`

	// PID is the path to the PID file.
	PID string `json:"pid"`

	// ReadTimeout is the read timeout for the server.
	ReadTimeout xtime.Duration `json:"readTimeout"`

	// WriteTimeout is the write timeout for the server.
	WriteTimeout xtime.Duration `json:"writeTimeout"`

	// IdleTimeout is the idle timeout for the server.
	IdleTimeout xtime.Duration `json:"idleTimeout"`
}

// validate validates the server configuration.
func (s *Server) validate() error {
	if s.TLS != nil {
		if s.TLS.Certificate == nil {
			return ErrServerMissingTLSCertificate
		}

		if s.TLS.Key == nil {
			return ErrServerMissingTLSKey
		}
	}

	if s.TLS == nil {
		return ErrServerMissingTLSConfig
	}

	return nil
}

// setSecrets populates the configuration with the necessary secrets.
func (s *Server) setSecrets(store credential.Store) error {
	certificate, err := store.GetBytes("tlscertificate")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	key, err := store.GetBytes("tlskey")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	s.TLS = &TLS{
		Certificate: certificate,
		Key:         key,
	}

	return nil
}

// setDefaults sets the default values for the server configuration.
func (s *Server) setDefaults() {
	if s.Address == "" {
		s.Address = defaultServerAddress
	}

	if s.PID == "" {
		s.PID = defaultServerPID
	}

	if s.ReadTimeout == 0 {
		s.ReadTimeout = defaultServerReadTimeout
	}

	if s.WriteTimeout == 0 {
		s.WriteTimeout = defaultServerWriteTimeout
	}

	if s.IdleTimeout == 0 {
		s.IdleTimeout = defaultServerIdleTimeout
	}
}
