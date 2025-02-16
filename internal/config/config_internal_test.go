// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config

import (
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xtime"
	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
	"go.cipher.host/pkgdex/internal/testhelper"
)

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    *Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			give: &Config{
				Service: &Service{
					Name:          "Example",
					Homepage:      "https://example.com",
					Contact:       "support@example.com",
					PrivacyPolicy: "https://example.com/privacy",
					Key:           "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				},
				Server: &Server{
					TLS: &TLS{
						Certificate: []byte("cert_data"),
						Key:         []byte("key_data"),
					},
					Address:      ":8080",
					PID:          "/var/run/pkgdex.pid",
					ReadTimeout:  xtime.Duration(5 * time.Second),
					WriteTimeout: xtime.Duration(10 * time.Second),
					IdleTimeout:  xtime.Duration(30 * time.Second),
				},
				Database: &Database{
					Path: "/var/lib/pkgdex/pkgdex.db",
				},
				Packages: []*Package{
					{
						Name:        "example",
						Description: "Example package",
						Version:     "1.0.0",
						Branch:      "trunk",
						Repository:  "https://git.example.com/example.git",
						License:     "MIT",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "missing service configuration",
			give:    &Config{},
			wantErr: true,
		},
		{
			name: "invalid service configuration",
			give: &Config{
				Service: &Service{
					Contact: "https://example.com/contact",
				},
			},
			wantErr: true,
		},
		{
			name: "missing server configuration",
			give: &Config{
				Service: &Service{
					Contact: "example@example.com",
					Key:     "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid server configuration",
			give: &Config{
				Service: &Service{
					Contact: "example@example.com",
					Key:     "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				},
				Server: &Server{},
			},
			wantErr: true,
		},
		{
			name: "missing database configuration",
			give: &Config{
				Service: &Service{
					Contact: "example@example.com",
					Key:     "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				},
				Server: &Server{
					TLS: &TLS{
						Certificate: []byte("data"),
						Key:         []byte("data"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid database configuration",
			give: &Config{
				Service: &Service{
					Contact: "example@example.com",
					Key:     "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				},
				Server: &Server{
					TLS: &TLS{
						Certificate: []byte("data"),
						Key:         []byte("data"),
					},
				},
				Database: &Database{},
			},
			wantErr: true,
		},
		{
			name: "missing packages configuration",
			give: &Config{
				Service: &Service{
					Contact: "example@example.com",
					Key:     "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				},
				Server: &Server{
					TLS: &TLS{
						Certificate: []byte("data"),
						Key:         []byte("data"),
					},
				},
				Database: &Database{
					Path: "/var/lib/pkgdex/pkgdex.db",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid packages configuration",
			give: &Config{
				Service: &Service{
					Contact: "example@example.com",
					Key:     "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				},
				Server: &Server{
					TLS: &TLS{
						Certificate: []byte("data"),
						Key:         []byte("data"),
					},
				},
				Database: &Database{
					Path: "/var/lib/pkgdex/pkgdex.db",
				},
				Packages: []*Package{
					{
						Description: "Example package",
						Repository:  "https://git.example.com/example.git",
						License:     "MIT",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.give.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_SetSecrets(t *testing.T) {
	testConfig := &Config{
		Service: &Service{},
		Server: &Server{
			Address: defaultServerAddress,
		},
		Database: &Database{
			Path: "/var/lib/pkgdex/pkgdex.db",
		},
	}

	tests := []struct {
		name    string
		give    *Config
		setup   func(t *testing.T, dir string) error
		want    *Config
		wantErr bool
		skipEnv bool
	}{
		{
			name: "set secrets successfully",
			give: testConfig,
			setup: func(t *testing.T, dir string) error {
				t.Helper()

				return testhelper.SetupSecrets(t, dir, map[string]string{
					"tlscertificate": "cert_data",
					"tlskey":         "key_data",
					"key":            "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				})
			},
			want: &Config{
				Service: &Service{
					Key: "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				},
				Server: &Server{
					Address: defaultServerAddress,
					TLS: &TLS{
						Certificate: xunsafe.StringToBytes("cert_data"),
						Key:         xunsafe.StringToBytes("key_data"),
					},
				},
				Database: &Database{
					Path: "/var/lib/pkgdex/pkgdex.db",
				},
			},
			wantErr: false,
		},
		{
			name: "fail to set TLS certificate",
			give: testConfig,
			setup: func(t *testing.T, dir string) error {
				t.Helper()

				return testhelper.SetupSecrets(t, dir, map[string]string{
					"tlskey": "key_data",
				})
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail to set credentials directory",
			give: testConfig,
			setup: func(t *testing.T, _ string) error {
				t.Helper()

				return nil
			},
			want:    nil,
			wantErr: true,
			skipEnv: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if err := tt.setup(t, tmpDir); err != nil {
				t.Fatal(err)
			}

			if !tt.skipEnv {
				t.Setenv("CREDENTIALS_DIRECTORY", tmpDir)
			}

			err := tt.give.setSecrets()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Config.SetSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.give, tt.want) {
				t.Fatalf("Config.SetSecrets() = %v, want %v", tt.give, tt.want)
			}
		})
	}
}
