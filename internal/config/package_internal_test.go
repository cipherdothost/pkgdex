// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config

import "testing"

func TestPackage_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    *Package
		wantErr bool
	}{
		{
			name: "valid package configuration",
			give: &Package{
				Name:        "example",
				Description: "This is an example package.",
				Version:     "1.0.0",
				Branch:      "trunk",
				Repository:  "https://git.sr.ht/~jamesponddotco/example",
				License:     "MIT",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			give: &Package{
				Description: "This is an example package.",
				Version:     "1.0.0",
				Branch:      "trunk",
				Repository:  "https://git.sr.ht/~jamesponddotco/example",
				License:     "MIT",
			},
			wantErr: true,
		},
		{
			name: "missing description",
			give: &Package{
				Name:       "example",
				Version:    "1.0.0",
				Branch:     "trunk",
				Repository: "https://git.sr.ht/~jamesponddotco/example",
				License:    "MIT",
			},
			wantErr: true,
		},
		{
			name: "missing version",
			give: &Package{
				Name:        "example",
				Description: "This is an example package.",
				Branch:      "trunk",
				Repository:  "https://git.sr.ht/~jamesponddotco/example",
				License:     "MIT",
			},
			wantErr: true,
		},
		{
			name: "missing branch",
			give: &Package{
				Name:        "example",
				Description: "This is an example package.",
				Version:     "1.0.0",
				Repository:  "https://git.sr.ht/~jamesponddotco/example",
				License:     "MIT",
			},
			wantErr: true,
		},
		{
			name: "missing repository",
			give: &Package{
				Name:        "example",
				Description: "This is an example package.",
				Version:     "1.0.0",
				Branch:      "trunk",
				License:     "MIT",
			},
			wantErr: true,
		},
		{
			name: "missing license",
			give: &Package{
				Name:        "example",
				Description: "This is an example package.",
				Version:     "1.0.0",
				Branch:      "trunk",
				Repository:  "https://git.sr.ht/~jamesponddotco/example",
			},
			wantErr: true,
		},
		{
			name: "invalid name",
			give: &Package{
				Name:        "example!",
				Description: "This is an example package.",
				Version:     "1.0.0",
				Branch:      "trunk",
				Repository:  "https://git.sr.ht/~jamesponddotco/example",
				License:     "MIT",
			},
			wantErr: true,
		},
		{
			name: "invalid version",
			give: &Package{
				Name:        "example",
				Description: "This is an example package.",
				Version:     "v1.0.0",
				Branch:      "trunk",
				Repository:  "https://git.sr.ht/~jamesponddotco/example",
				License:     "MIT",
			},
			wantErr: true,
		},
		{
			name: "invalid repository URL",
			give: &Package{
				Name:        "example",
				Description: "This is an example package.",
				Version:     "1.0.0",
				Branch:      "trunk",
				Repository:  "://git.sr.ht/~jamesponddotco/example",
				License:     "MIT",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.give.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Package.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
