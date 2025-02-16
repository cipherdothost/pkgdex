// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"git.sr.ht/~jamesponddotco/credential-go"
	"go.cipher.host/pkgdex/internal/meta"
)

func TestServer_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    *Server
		wantErr bool
	}{
		{
			name: "valid server configuration",
			give: &Server{
				TLS: &TLS{
					Certificate: []byte("data"),
					Key:         []byte("data"),
				},
				Address:      defaultServerAddress,
				PID:          defaultServerPID,
				ReadTimeout:  defaultServerReadTimeout,
				WriteTimeout: defaultServerWriteTimeout,
				IdleTimeout:  defaultServerIdleTimeout,
			},
			wantErr: false,
		},
		{
			name: "missing TLS configuration",
			give: &Server{
				Address:      defaultServerAddress,
				PID:          defaultServerPID,
				ReadTimeout:  defaultServerReadTimeout,
				WriteTimeout: defaultServerWriteTimeout,
				IdleTimeout:  defaultServerIdleTimeout,
			},
			wantErr: true,
		},
		{
			name: "missing TLS certificate",
			give: &Server{
				TLS: &TLS{
					Key: []byte("data"),
				},
				Address:      defaultServerAddress,
				PID:          defaultServerPID,
				ReadTimeout:  defaultServerReadTimeout,
				WriteTimeout: defaultServerWriteTimeout,
				IdleTimeout:  defaultServerIdleTimeout,
			},
			wantErr: true,
		},
		{
			name: "missing TLS key",
			give: &Server{
				TLS: &TLS{
					Certificate: []byte("data"),
				},
				Address:      defaultServerAddress,
				PID:          defaultServerPID,
				ReadTimeout:  defaultServerReadTimeout,
				WriteTimeout: defaultServerWriteTimeout,
				IdleTimeout:  defaultServerIdleTimeout,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.give.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Server.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_SetSecrets(t *testing.T) {
	tests := []struct {
		name    string
		give    *Server
		setup   func(dir string) error
		want    *Server
		wantErr bool
	}{
		{
			name: "set secrets successfully",
			give: &Server{
				Address:      defaultServerAddress,
				PID:          defaultServerPID,
				ReadTimeout:  defaultServerReadTimeout,
				WriteTimeout: defaultServerWriteTimeout,
				IdleTimeout:  defaultServerIdleTimeout,
			},
			setup: func(dir string) error {
				if err := os.WriteFile(filepath.Join(dir, "pkgdex-tlscertificate"), []byte("data"), 0o600); err != nil {
					return err
				}

				return os.WriteFile(filepath.Join(dir, "pkgdex-tlskey"), []byte("data"), 0o600)
			},
			want: &Server{
				TLS: &TLS{
					Certificate: []byte("data"),
					Key:         []byte("data"),
				},
				Address:      defaultServerAddress,
				PID:          defaultServerPID,
				ReadTimeout:  defaultServerReadTimeout,
				WriteTimeout: defaultServerWriteTimeout,
				IdleTimeout:  defaultServerIdleTimeout,
			},
			wantErr: false,
		},
		{
			name: "fail to set TLS certificate",
			give: &Server{
				Address:      defaultServerAddress,
				PID:          defaultServerPID,
				ReadTimeout:  defaultServerReadTimeout,
				WriteTimeout: defaultServerWriteTimeout,
				IdleTimeout:  defaultServerIdleTimeout,
			},
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "pkgdex-tlskey"), []byte("data"), 0o600)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail to set TLS key",
			give: &Server{
				TLS: &TLS{
					Certificate: []byte("data"),
					Key:         []byte("data"),
				},
				Address:      defaultServerAddress,
				PID:          defaultServerPID,
				ReadTimeout:  defaultServerReadTimeout,
				WriteTimeout: defaultServerWriteTimeout,
				IdleTimeout:  defaultServerIdleTimeout,
			},
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "pkgdex-tlscertificate"), []byte("data"), 0o600)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if err := tt.setup(tmpDir); err != nil {
				t.Fatal(err)
			}

			t.Setenv("CREDENTIALS_DIRECTORY", tmpDir)

			dir, err := credential.Open(meta.Name)
			if err != nil {
				t.Fatal(err)
			}

			err = tt.give.setSecrets(dir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Server.SetSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.give, tt.want) {
				t.Fatalf("Server.SetSecrets() = %v, want %v", tt.give, tt.want)
			}
		})
	}
}

func TestServer_SetDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give *Server
		want *Server
	}{
		{
			name: "set defaults successfully",
			give: &Server{},
			want: &Server{
				Address:      defaultServerAddress,
				PID:          defaultServerPID,
				ReadTimeout:  defaultServerReadTimeout,
				WriteTimeout: defaultServerWriteTimeout,
				IdleTimeout:  defaultServerIdleTimeout,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.give.setDefaults()

			if !reflect.DeepEqual(tt.give, tt.want) {
				t.Fatalf("Server.SetDefaults() = %v, want %v", tt.give, tt.want)
			}
		})
	}
}
