// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package repository_test

import (
	"errors"
	"testing"

	"go.cipher.host/pkgdex/internal/repository"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		giveURL    string
		giveBranch string
		want       *repository.Repository
		wantErr    error
	}{
		{
			name:       "valid GitHub URL with HTTPS",
			giveURL:    "https://github.com/user/repo",
			giveBranch: "main",
			want: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			wantErr: nil,
		},
		{
			name:       "valid GitHub URL without protocol",
			giveURL:    "github.com/user/repo",
			giveBranch: "main",
			want: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			wantErr: nil,
		},
		{
			name:       "valid Sourcehut URL with HTTPS",
			giveURL:    "https://git.sr.ht/~user/repo",
			giveBranch: "trunk",
			want: &repository.Repository{
				Service: repository.ServiceSourcehut,
				Owner:   "~user",
				Name:    "repo",
				Branch:  "trunk",
			},
			wantErr: nil,
		},
		{
			name:       "valid Sourcehut URL without protocol",
			giveURL:    "git.sr.ht/~user/repo",
			giveBranch: "trunk",
			want: &repository.Repository{
				Service: repository.ServiceSourcehut,
				Owner:   "~user",
				Name:    "repo",
				Branch:  "trunk",
			},
			wantErr: nil,
		},
		{
			name:       "valid GitLab URL with HTTPS",
			giveURL:    "https://gitlab.com/user/repo",
			giveBranch: "master",
			want: &repository.Repository{
				Service: repository.ServiceGitLab,
				Owner:   "user",
				Name:    "repo",
				Branch:  "master",
			},
			wantErr: nil,
		},
		{
			name:       "valid GitLab URL without protocol",
			giveURL:    "gitlab.com/user/repo",
			giveBranch: "master",
			want: &repository.Repository{
				Service: repository.ServiceGitLab,
				Owner:   "user",
				Name:    "repo",
				Branch:  "master",
			},
			wantErr: nil,
		},
		{
			name:       "URL with www prefix",
			giveURL:    "https://www.github.com/user/repo",
			giveBranch: "main",
			want: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			wantErr: nil,
		},
		{
			name:       "URL with trailing slash",
			giveURL:    "https://github.com/user/repo/",
			giveBranch: "main",
			want: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			wantErr: nil,
		},
		{
			name:       "invalid URL",
			giveURL:    "://invalid.com/user/repo",
			giveBranch: "main",
			want:       nil,
			wantErr:    repository.ErrInvalidURL,
		},
		{
			name:       "URL without owner and name",
			giveURL:    "https://github.com",
			giveBranch: "main",
			want:       nil,
			wantErr:    repository.ErrInvalidURL,
		},
		{
			name:       "URL with owner but no name",
			giveURL:    "https://github.com/user",
			giveBranch: "main",
			want:       nil,
			wantErr:    repository.ErrInvalidURL,
		},
		{
			name:       "URL with extra path segments",
			giveURL:    "https://github.com/user/repo/extra/path",
			giveBranch: "main",
			want: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			wantErr: nil,
		},
		{
			name:       "empty URL",
			giveURL:    "",
			giveBranch: "main",
			want:       nil,
			wantErr:    repository.ErrMissingURL,
		},
		{
			name:       "empty branch",
			giveURL:    "https://github.com/user/repo",
			giveBranch: "",
			want:       nil,
			wantErr:    repository.ErrMissingBranch,
		},
		{
			name:       "URL with query parameters",
			giveURL:    "https://github.com/user/repo?param=value",
			giveBranch: "main",
			want: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			wantErr: nil,
		},
		{
			name:       "URL with fragment",
			giveURL:    "https://github.com/user/repo#fragment",
			giveBranch: "main",
			want: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			wantErr: nil,
		},
		{
			name:       "unsupported service",
			giveURL:    "https://unsupported.com/user/repo",
			giveBranch: "main",
			want:       nil,
			wantErr:    repository.ErrUnsupportedService,
		},
		{
			name:       "unsupported service without protocol",
			giveURL:    "unsupported.com/user/repo",
			giveBranch: "main",
			want:       nil,
			wantErr:    repository.ErrUnsupportedService,
		},
		{
			name:       "unsupported service with www prefix",
			giveURL:    "https://www.unsupported.com/user/repo",
			giveBranch: "main",
			want:       nil,
			wantErr:    repository.ErrUnsupportedService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := repository.Parse(tt.giveURL, tt.giveBranch)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Parse(%q, %q) error = %v, wantErr %v", tt.giveURL, tt.giveBranch, err, tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Errorf("Parse(%q, %q) unexpected error: %v", tt.giveURL, tt.giveBranch, err)
			}

			if got.Service != tt.want.Service {
				t.Errorf("Parse(%q, %q) Service = %v, want %v", tt.giveURL, tt.giveBranch, got.Service, tt.want.Service)
			}

			if got.Owner != tt.want.Owner {
				t.Errorf("Parse(%q, %q) Owner = %v, want %v", tt.giveURL, tt.giveBranch, got.Owner, tt.want.Owner)
			}

			if got.Name != tt.want.Name {
				t.Errorf("Parse(%q, %q) Name = %v, want %v", tt.giveURL, tt.giveBranch, got.Name, tt.want.Name)
			}

			if got.Branch != tt.want.Branch {
				t.Errorf("Parse(%q, %q) Branch = %v, want %v", tt.giveURL, tt.giveBranch, got.Branch, tt.want.Branch)
			}
		})
	}
}

func TestRepository_CloneURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give *repository.Repository
		want string
	}{
		{
			name: "GitHub repository",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			want: "https://github.com/user/repo",
		},
		{
			name: "Sourcehut repository",
			give: &repository.Repository{
				Service: repository.ServiceSourcehut,
				Owner:   "~user",
				Name:    "repo",
				Branch:  "trunk",
			},
			want: "https://git.sr.ht/~user/repo",
		},
		{
			name: "GitLab repository",
			give: &repository.Repository{
				Service: repository.ServiceGitLab,
				Owner:   "user",
				Name:    "repo",
				Branch:  "master",
			},
			want: "https://gitlab.com/user/repo",
		},
		{
			name: "GitHub repository with organization",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "organization",
				Name:    "repo",
				Branch:  "main",
			},
			want: "https://github.com/organization/repo",
		},
		{
			name: "repository with dashes",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user-name",
				Name:    "repo-name",
				Branch:  "main",
			},
			want: "https://github.com/user-name/repo-name",
		},
		{
			name: "repository with dots",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user.name",
				Name:    "repo.name",
				Branch:  "main",
			},
			want: "https://github.com/user.name/repo.name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.give.CloneURL()
			if got != tt.want {
				t.Errorf("CloneURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_SourceURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		give     *repository.Repository
		giveFile string
		giveLine string
		want     string
	}{
		{
			name: "GitHub without line number",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			giveFile: "cmd/app/main.go",
			giveLine: "",
			want:     "https://github.com/user/repo/blob/main/cmd/app/main.go",
		},
		{
			name: "GitHub with line number",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			giveFile: "cmd/app/main.go",
			giveLine: "42",
			want:     "https://github.com/user/repo/blob/main/cmd/app/main.go#L42",
		},
		{
			name: "Sourcehut without line number",
			give: &repository.Repository{
				Service: repository.ServiceSourcehut,
				Owner:   "~user",
				Name:    "repo",
				Branch:  "trunk",
			},
			giveFile: "cmd/app/main.go",
			giveLine: "",
			want:     "https://git.sr.ht/~user/repo/tree/trunk/item/cmd/app/main.go",
		},
		{
			name: "Sourcehut with line number",
			give: &repository.Repository{
				Service: repository.ServiceSourcehut,
				Owner:   "~user",
				Name:    "repo",
				Branch:  "trunk",
			},
			giveFile: "cmd/app/main.go",
			giveLine: "42",
			want:     "https://git.sr.ht/~user/repo/tree/trunk/item/cmd/app/main.go#L42",
		},
		{
			name: "GitLab without line number",
			give: &repository.Repository{
				Service: repository.ServiceGitLab,
				Owner:   "user",
				Name:    "repo",
				Branch:  "master",
			},
			giveFile: "cmd/app/main.go",
			giveLine: "",
			want:     "https://gitlab.com/user/repo/-/blob/master/cmd/app/main.go",
		},
		{
			name: "GitLab with line number",
			give: &repository.Repository{
				Service: repository.ServiceGitLab,
				Owner:   "user",
				Name:    "repo",
				Branch:  "master",
			},
			giveFile: "cmd/app/main.go",
			giveLine: "42",
			want:     "https://gitlab.com/user/repo/-/blob/master/cmd/app/main.go#L42",
		},
		{
			name: "file with special characters",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			giveFile: "path/with spaces/and#symbols.go",
			giveLine: "42",
			want:     "https://github.com/user/repo/blob/main/path/with spaces/and#symbols.go#L42",
		},
		{
			name: "empty file path",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			giveFile: "",
			giveLine: "42",
			want:     "https://github.com/user/repo/blob/main/#L42",
		},
		{
			name: "unsupported service",
			give: &repository.Repository{
				Service: "unsupported.com",
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			giveFile: "main.go",
			giveLine: "42",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.give.SourceURL(tt.giveFile, tt.giveLine)
			if got != tt.want {
				t.Errorf("SourceURL(%q, %q) = %q, want %q", tt.giveFile, tt.giveLine, got, tt.want)
			}
		})
	}
}

func TestRepository_GoImportMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		give       *repository.Repository
		giveImport string
		want       string
	}{
		{
			name: "GitHub repository",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			giveImport: "example.com/user/repo",
			want:       "example.com/user/repo git https://github.com/user/repo",
		},
		{
			name: "Sourcehut repository",
			give: &repository.Repository{
				Service: repository.ServiceSourcehut,
				Owner:   "~user",
				Name:    "repo",
				Branch:  "trunk",
			},
			giveImport: "example.com/user/repo",
			want:       "example.com/user/repo git https://git.sr.ht/~user/repo",
		},
		{
			name: "GitLab repository",
			give: &repository.Repository{
				Service: repository.ServiceGitLab,
				Owner:   "user",
				Name:    "repo",
				Branch:  "master",
			},
			giveImport: "example.com/user/repo",
			want:       "example.com/user/repo git https://gitlab.com/user/repo",
		},
		{
			name: "with complex import path",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			giveImport: "custom.domain.com/path/to/repo",
			want:       "custom.domain.com/path/to/repo git https://github.com/user/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.give.GoImportMeta(tt.giveImport)
			if got != tt.want {
				t.Errorf("GoImportMeta(%q) = %q, want %q",
					tt.giveImport, got, tt.want)
			}
		})
	}
}

func TestRepository_GoSourceMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		give       *repository.Repository
		giveImport string
		want       string
	}{
		{
			name: "GitHub repository",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			giveImport: "example.com/user/repo",
			want:       "example.com/user/repo https://github.com/user/repo https://github.com/user/repo/tree/main{/dir} https://github.com/user/repo/blob/main{/dir}/{file}#L{line}",
		},
		{
			name: "Sourcehut repository",
			give: &repository.Repository{
				Service: repository.ServiceSourcehut,
				Owner:   "~user",
				Name:    "repo",
				Branch:  "trunk",
			},
			giveImport: "example.com/user/repo",
			want:       "example.com/user/repo https://git.sr.ht/~user/repo https://git.sr.ht/~user/repo/tree/trunk/item{/dir} https://git.sr.ht/~user/repo/tree/trunk/item{/dir}/{file}#L{line}",
		},
		{
			name: "GitLab repository",
			give: &repository.Repository{
				Service: repository.ServiceGitLab,
				Owner:   "user",
				Name:    "repo",
				Branch:  "master",
			},
			giveImport: "example.com/user/repo",
			want:       "example.com/user/repo https://gitlab.com/user/repo https://gitlab.com/user/repo/-/tree/master{/dir} https://gitlab.com/user/repo/-/blob/master{/dir}/{file}#L{line}",
		},
		{
			name: "with complex import path",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "main",
			},
			giveImport: "custom.domain.com/path/to/repo",
			want:       "custom.domain.com/path/to/repo https://github.com/user/repo https://github.com/user/repo/tree/main{/dir} https://github.com/user/repo/blob/main{/dir}/{file}#L{line}",
		},
		{
			name: "with organization owner",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "organization",
				Name:    "repo",
				Branch:  "main",
			},
			giveImport: "example.com/organization/repo",
			want:       "example.com/organization/repo https://github.com/organization/repo https://github.com/organization/repo/tree/main{/dir} https://github.com/organization/repo/blob/main{/dir}/{file}#L{line}",
		},
		{
			name: "with development branch",
			give: &repository.Repository{
				Service: repository.ServiceGitHub,
				Owner:   "user",
				Name:    "repo",
				Branch:  "development",
			},
			giveImport: "example.com/user/repo",
			want:       "example.com/user/repo https://github.com/user/repo https://github.com/user/repo/tree/development{/dir} https://github.com/user/repo/blob/development{/dir}/{file}#L{line}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.give.GoSourceMeta(tt.giveImport)
			if got != tt.want {
				t.Errorf("GoSourceMeta(%q) = %q, want %q",
					tt.giveImport, got, tt.want)
			}
		})
	}
}
