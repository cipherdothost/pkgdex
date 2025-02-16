// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package pagination_test

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"go.cipher.host/pkgdex/internal/pagination"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		giveTotalItems   int
		giveCurrentPage  int
		giveItemsPerPage int
		want             pagination.Info
	}{
		{
			name:             "normal case",
			giveTotalItems:   100,
			giveCurrentPage:  3,
			giveItemsPerPage: 10,
			want: pagination.Info{
				Pages: []int{
					1,
					2,
					3,
					4,
					5,
				},
				CurrentPage:  3,
				TotalPages:   10,
				ItemsPerPage: 10,
				HasPrev:      true,
				HasNext:      true,
			},
		},
		{
			name:             "first page",
			giveTotalItems:   100,
			giveCurrentPage:  1,
			giveItemsPerPage: 10,
			want: pagination.Info{
				Pages: []int{
					1,
					2,
					3,
					4,
					5,
				},
				CurrentPage:  1,
				TotalPages:   10,
				ItemsPerPage: 10,
				HasPrev:      false,
				HasNext:      true,
			},
		},
		{
			name:             "last page",
			giveTotalItems:   100,
			giveCurrentPage:  10,
			giveItemsPerPage: 10,
			want: pagination.Info{
				Pages: []int{
					6,
					7,
					8,
					9,
					10,
				},
				CurrentPage:  10,
				TotalPages:   10,
				ItemsPerPage: 10,
				HasPrev:      true,
				HasNext:      false,
			},
		},
		{
			name:             "current page too high",
			giveTotalItems:   100,
			giveCurrentPage:  15,
			giveItemsPerPage: 10,
			want: pagination.Info{
				Pages: []int{
					6,
					7,
					8,
					9,
					10,
				},
				CurrentPage:  10,
				TotalPages:   10,
				ItemsPerPage: 10,
				HasPrev:      true,
				HasNext:      false,
			},
		},
		{
			name:             "current page too low",
			giveTotalItems:   100,
			giveCurrentPage:  0,
			giveItemsPerPage: 10,
			want: pagination.Info{
				Pages: []int{
					1,
					2,
					3,
					4,
					5,
				},
				CurrentPage:  1,
				TotalPages:   10,
				ItemsPerPage: 10,
				HasPrev:      false,
				HasNext:      true,
			},
		},
		{
			name:             "no items",
			giveTotalItems:   0,
			giveCurrentPage:  1,
			giveItemsPerPage: 10,
			want: pagination.Info{
				Pages: []int{
					1,
				},
				CurrentPage:  1,
				TotalPages:   1,
				ItemsPerPage: 10,
				HasPrev:      false,
				HasNext:      false,
			},
		},
		{
			name:             "single page",
			giveTotalItems:   5,
			giveCurrentPage:  1,
			giveItemsPerPage: 10,
			want: pagination.Info{
				Pages: []int{
					1,
				},
				CurrentPage:  1,
				TotalPages:   1,
				ItemsPerPage: 10,
				HasPrev:      false,
				HasNext:      false,
			},
		},
		{
			name:             "exact division",
			giveTotalItems:   50,
			giveCurrentPage:  3,
			giveItemsPerPage: 10,
			want: pagination.Info{
				Pages: []int{
					1,
					2,
					3,
					4,
					5,
				},
				CurrentPage:  3,
				TotalPages:   5,
				ItemsPerPage: 10,
				HasPrev:      true,
				HasNext:      true,
			},
		},
		{
			name:             "inexact division",
			giveTotalItems:   55,
			giveCurrentPage:  3,
			giveItemsPerPage: 10,
			want: pagination.Info{
				Pages: []int{
					1,
					2,
					3,
					4,
					5,
				},
				CurrentPage:  3,
				TotalPages:   6,
				ItemsPerPage: 10,
				HasPrev:      true,
				HasNext:      true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := pagination.New(tt.giveTotalItems, tt.giveCurrentPage, tt.giveItemsPerPage)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestInfo_Offset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		info pagination.Info
		want int
	}{
		{
			name: "first page",
			info: pagination.Info{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			want: 0,
		},
		{
			name: "second page",
			info: pagination.Info{
				CurrentPage:  2,
				ItemsPerPage: 10,
			},
			want: 10,
		},
		{
			name: "third page",
			info: pagination.Info{
				CurrentPage:  3,
				ItemsPerPage: 15,
			},
			want: 30,
		},
		{
			name: "custom items per page",
			info: pagination.Info{
				CurrentPage:  2,
				ItemsPerPage: 25,
			},
			want: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.info.Offset(); got != tt.want {
				t.Errorf("Info.Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_Limit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		info pagination.Info
		want int
	}{
		{
			name: "standard limit",
			info: pagination.Info{
				ItemsPerPage: 10,
			},
			want: 10,
		},
		{
			name: "custom limit",
			info: pagination.Info{
				ItemsPerPage: 25,
			},
			want: 25,
		},
		{
			name: "zero limit",
			info: pagination.Info{
				ItemsPerPage: 0,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.info.Limit(); got != tt.want {
				t.Errorf("Info.Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		urlStr  string
		want    int
		wantErr bool
	}{
		{
			name:   "valid page",
			urlStr: "http://example.com?page=5",
			want:   5,
		},
		{
			name:   "no page parameter",
			urlStr: "http://example.com",
			want:   1,
		},
		{
			name:   "empty page parameter",
			urlStr: "http://example.com?page=",
			want:   1,
		},
		{
			name:   "invalid page number",
			urlStr: "http://example.com?page=invalid",
			want:   1,
		},
		{
			name:   "negative page number",
			urlStr: "http://example.com?page=-1",
			want:   1,
		},
		{
			name:   "zero page number",
			urlStr: "http://example.com?page=0",
			want:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uri, err := url.Parse(tt.urlStr)
			if err != nil {
				t.Fatalf("GetPage() failed to parse URL: %v", err)
			}

			req := &http.Request{
				URL: uri,
			}
			if got := pagination.GetPage(req); got != tt.want {
				t.Errorf("GetPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
