// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config_test

import (
	"reflect"
	"testing"
	"time"

	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/testhelper"
	"go.cipher.host/x/xtime"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		give    string
		setup   func(t *testing.T, dir string) error
		want    *config.Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			give: "testdata/valid-config.json",
			setup: func(t *testing.T, dir string) error {
				t.Helper()

				return testhelper.SetupSecrets(t, dir, map[string]string{
					"tlscertificate": "cert_data",
					"tlskey":         "key_data",
					"key":            "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
				})
			},
			want: &config.Config{
				Service: &config.Service{
					Name:            "Example",
					Description:     "A centralized service providing custom import paths and essential metadata for Go packages.",
					Homepage:        "https://example.com",
					BaseURL:         "example.com",
					Contact:         "support@example.com",
					PrivacyPolicy:   "https://example.com/privacy",
					Key:             "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9",
					PackagesPerPage: 2,
				},
				Server: &config.Server{
					TLS: &config.TLS{
						Certificate: []byte("cert_data"),
						Key:         []byte("key_data"),
					},
					Address:      ":8080",
					PID:          "/var/run/pkgdex.pid",
					ReadTimeout:  xtime.Duration(5 * time.Second),
					WriteTimeout: xtime.Duration(10 * time.Second),
					IdleTimeout:  xtime.Duration(30 * time.Second),
				},
				Database: &config.Database{
					Path:      "/var/lib/pkgdex/pkgdex.db",
					IndexPath: "/var/lib/pkgdex/index",
				},
				Packages: []*config.Package{
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
			name: "invalid configuration",
			give: "testdata/invalid-config.json",
			setup: func(t *testing.T, dir string) error {
				t.Helper()

				return testhelper.SetupSecrets(t, dir, map[string]string{
					"tlscertificate": "cert_data",
					"tlskey":         "key_data",
				})
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non-existent file",
			give: "testdata/non-existent.json",
			setup: func(_ *testing.T, _ string) error {
				return nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid JSON file",
			give: "testdata/invalid-json.json",
			setup: func(_ *testing.T, _ string) error {
				return nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing required fields",
			give: "testdata/missing-fields.json",
			setup: func(_ *testing.T, _ string) error {
				return nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing secrets",
			give: "testdata/valid-config.json",
			setup: func(_ *testing.T, _ string) error {
				return nil
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			t.Setenv("CREDENTIALS_DIRECTORY", tmpDir)

			if err := tt.setup(t, tmpDir); err != nil {
				t.Fatalf("Failed to set up test: %v", err)
			}

			got, err := config.Load(tt.give)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Load() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
