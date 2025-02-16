// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package testhelper_test

import (
	"testing"

	"go.cipher.host/pkgdex/internal/testhelper"
)

func TestSetupSecrets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		giveDirectory string
		giveSecrets   map[string]string
		wantError     bool
	}{
		{
			name:          "valid without secrets",
			giveDirectory: t.TempDir(),
			giveSecrets:   make(map[string]string),
			wantError:     false,
		},
		{
			name:          "valid with secrets",
			giveDirectory: t.TempDir(),
			giveSecrets: map[string]string{
				"SECRET":         "secret",
				"ANOTHER_SECRET": "another secret",
			},
			wantError: false,
		},
		{
			name:          "invalid directory",
			giveDirectory: "invalid",
			giveSecrets: map[string]string{
				"SECRET": "secret",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := testhelper.SetupSecrets(t, tt.giveDirectory, tt.giveSecrets)
			if (err != nil) != tt.wantError {
				t.Fatalf("SetupSecrets() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}
