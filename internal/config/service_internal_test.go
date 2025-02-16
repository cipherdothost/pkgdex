// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config

import (
	"reflect"
	"testing"
)

func TestService_Validate(t *testing.T) {
	t.Parallel()

	defaultServiceKey := "pkgdex_tmWhkGPm3PB4e9yju8FVEhBV45v9JKSmPWRHxRU9gXFh4NUse8e3b13c9"

	tests := []struct {
		name    string
		give    *Service
		wantErr bool
	}{
		{
			name: "valid service configuration",
			give: &Service{
				Name:     defaultServiceName,
				Homepage: defaultServiceHomepage,
				Contact:  defaultServiceContact,
				Key:      defaultServiceKey,
			},
			wantErr: false,
		},
		{
			name: "invalid email address",
			give: &Service{
				Name:     defaultServiceName,
				Homepage: defaultServiceHomepage,
				Contact:  "invalid",
				Key:      defaultServiceKey,
			},
			wantErr: true,
		},
		{
			name: "invalid homepage URL",
			give: &Service{
				Name:     defaultServiceName,
				Homepage: "://invalid.com",
				Contact:  defaultServiceContact,
				Key:      defaultServiceKey,
			},
			wantErr: true,
		},
		{
			name: "invalid base URL",
			give: &Service{
				Name:     defaultServiceName,
				Homepage: defaultServiceHomepage,
				Contact:  defaultServiceContact,
				BaseURL:  "://invalid.com",
				Key:      defaultServiceKey,
			},
			wantErr: true,
		},
		{
			name: "invalid privacy policy URL",
			give: &Service{
				Name:          defaultServiceName,
				Homepage:      defaultServiceHomepage,
				Contact:       defaultServiceContact,
				PrivacyPolicy: "://invalid.com/privacy",
				Key:           defaultServiceKey,
			},
			wantErr: true,
		},
		{
			name: "missing service key",
			give: &Service{
				Name:     defaultServiceName,
				Homepage: defaultServiceHomepage,
				Contact:  defaultServiceContact,
			},
			wantErr: true,
		},
		{
			name: "small and invalid service key",
			give: &Service{
				Name:     defaultServiceName,
				Homepage: defaultServiceHomepage,
				Contact:  defaultServiceContact,
				Key:      "small",
			},
			wantErr: true,
		},
		{
			name: "correct size, but still invalid service key",
			give: &Service{
				Name:     defaultServiceName,
				Homepage: defaultServiceHomepage,
				Contact:  defaultServiceContact,
				Key:      "nDQyvamQbUTYrOmQH1FLIBuP4pkA4kPd0gRprNiDXSPRo4lXaLeQowV8ey4spGbp",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.give.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Service.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_SetDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give *Service
		want *Service
	}{
		{
			name: "set defaults successfully",
			give: &Service{},
			want: &Service{
				Name:            defaultServiceName,
				Description:     defaultServiceDescription,
				Homepage:        defaultServiceHomepage,
				BaseURL:         defaultServiceBaseURL,
				Contact:         defaultServiceContact,
				PackagesPerPage: defaultPackagesPerPage,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.give.setDefaults()

			if !reflect.DeepEqual(tt.give, tt.want) {
				t.Fatalf("Service.SetDefaults() = %v, want %v", tt.give, tt.want)
			}
		})
	}
}
