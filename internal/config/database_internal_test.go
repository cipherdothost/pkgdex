// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package config

import (
	"reflect"
	"testing"
)

func TestDatabase_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    *Database
		wantErr bool
	}{
		{
			name: "valid database configuration",
			give: &Database{
				Path: "/path/to/database.db",
			},
			wantErr: false,
		},
		{
			name:    "missing path",
			give:    &Database{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.give.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Database.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_SetDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give *Database
		want *Database
	}{
		{
			name: "set defaults successfully",
			give: &Database{},
			want: &Database{
				Path:      defaultDatabasePath,
				IndexPath: defaultDatabaseIndexPath,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.give.setDefaults()

			if !reflect.DeepEqual(tt.give, tt.want) {
				t.Fatalf("Database.SetDefaults() = %v, want %v", tt.give, tt.want)
			}
		})
	}
}
