# SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
#
# SPDX-License-Identifier: EUPL-1.2

.POSIX:
.SUFFIXES:

PKGNAME := pkgdex
PKGBIN  := $(PKGNAME)ctl

PREFIX     := /usr/local
BINDIR     := bin
SHAREDIR   := share
MANDIR     := ${SHAREDIR}/man
PKGDIR     := ./
CONFIGDIR  := config
DOCDIR     := docs
OUTDIR     := build
ASSETSDIR  := static/assets

GO        ?= go
GIT       ?= git
REUSE     ?= reuse
RM        ?= rm
INSTALL   ?= install
NPM       ?= npm
SASS      ?= ./node_modules/.bin/sass
STYLELINT ?= ./node_modules/.bin/stylelint
SCDOC     ?= scdoc

REQUIRED_TOOLS := $(GO) $(GIT) $(SASS) $(SCDOC) $(NPM)
GO_MIN_VERSION := 1.24

GOBUILD_FLAGS := -trimpath \
				 -buildmode=pie \
				 -mod=readonly \
				 -modcacherw
LDFLAGS       := -s -w
GOBUILD_OUT   := -o $(OUTDIR)/$(PKGBIN)
GOBUILD_OPTS  := $(GOBUILD_FLAGS) -ldflags '$(LDFLAGS)' $(GOBUILD_OUT)

all: check build build/doc

check: # Checks that all required tools for building the application are installed.
	$(foreach tool,$(REQUIRED_TOOLS),\
		$(if $(shell command -v $(tool)),,$(error "$(tool) not found in PATH")))
	@$(GO) version | awk 'NR==1 {if ($$3 < "go$(GO_MIN_VERSION)") exit 1}'

pre-commit: tidy fmt lint vulnerabilities test build clean # Runs all pre-commit checks.

commit: pre-commit # Commits the changes to the repository.
	$(GIT) commit -s

build: check build/css # Builds an application binary.
	$(GO) build $(GOBUILD_OPTS) $(PKGDIR)

build/linux: check build/css # Builds an application binary for Linux.
	GOOS=linux GOARCH=amd64 $(GO) build $(GOBUILD_OPTS) $(PKGDIR)

build/css: # Builds the CSS assets.
	$(SASS) --no-source-map --quiet --style 'compressed' \
		${ASSETSDIR}/scss/main.scss \
		${ASSETSDIR}/css/main.css
	$(SASS) --no-source-map --quiet --style 'compressed' \
		${ASSETSDIR}/scss/highlight.scss \
		${ASSETSDIR}/css/highlight.css

build/doc: # Builds the manpages.
	$(SCDOC) <$(DOCDIR)/$(PKGNAME).1.scd >$(OUTDIR)/$(PKGNAME).1
	$(SCDOC) <$(DOCDIR)/$(PKGNAME).5.scd >$(OUTDIR)/$(PKGNAME).5

install: # Installs the release binary and other assets.
	$(INSTALL) -d \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/ \
		$(DESTDIR)$(PREFIX)/$(SHAREDIR)/$(PKGNAME)/ \
		$(DESTDIR)$(PREFIX)/$(SHAREDIR)/$(PKGNAME)/apparmor/ \
		$(DESTDIR)$(PREFIX)/$(SHAREDIR)/$(PKGNAME)/nginx/ \
		$(DESTDIR)$(PREFIX)/$(SHAREDIR)/$(PKGNAME)/systemd/ \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/ \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man5/
	$(INSTALL) -pm 0755 $(OUTDIR)/$(PKGBIN) \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/
	$(INSTALL) -pm 0644 $(CONFIGDIR)/config.example.json \
		$(DESTDIR)$(PREFIX)/$(SHAREDIR)/$(PKGNAME)/
	$(INSTALL) -pm 0644 $(CONFIGDIR)/apparmor/usr.local.bin.$(PKGNAME) \
		$(DESTDIR)$(PREFIX)/$(SHAREDIR)/$(PKGNAME)/apparmor/
	$(INSTALL) -pm 0644 $(CONFIGDIR)/nginx/pkg.example.com \
		$(DESTDIR)$(PREFIX)/$(SHAREDIR)/$(PKGNAME)/nginx/
	$(INSTALL) -pm 0644 $(CONFIGDIR)/systemd/$(PKGNAME).service \
		$(DESTDIR)$(PREFIX)/$(SHAREDIR)/$(PKGNAME)/systemd/
	$(INSTALL) -pm 0644 $(OUTDIR)/$(PKGNAME).1 \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/
	$(INSTALL) -pm 0644 $(OUTDIR)/$(PKGNAME).5 \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man5/

development: build/css # Starts the server for development.
	 CREDENTIALS_DIRECTORY='.dev' $(GO) run ${PKGDIR} start \
	 					   --config '.dev/config.json' || exit 0

tidy: # Updates the go.mod file to ensure it matches the source code.
	$(GO) mod tidy

fmt: # Formats Go source files in this repository.
	$(GO) run mvdan.cc/gofumpt@latest -e -extra -w .

lint: # Runs golangci-lint using the config at the root of the repository.
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...

lint/styles: # Runs stylelint on the CSS and SASS assets.
	$(STYLELINT) --quiet-deprecation-warnings \
		--ignore-pattern '*breakpoints*' \
		$(ASSETSDIR)/scss/

lint/licenses: # Run reuse to ensure the project complies with the REUSE specification.
	$(REUSE) lint --lines

fix/styles: # Runs stylelint on the CSS and SASS assets and fixes any issues.
	$(STYLELINT) --quiet-deprecation-warnings \
		--ignore-pattern '*breakpoints*' \
		--fix 'strict' \
		$(ASSETSDIR)/scss/

vulnerabilities: # Analyzes the codebase and looks for vulnerabilities affecting it.
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

test: # Runs unit tests.
	$(GO) test -cover -race -vet all -mod readonly ./...

test/coverage: # Generates a coverage profile and open it in a browser.
	$(GO) test -coverprofile cover.out ./...
	$(GO) tool cover -html=cover.out

licenses: # Runs go-licenses to check the licenses of the dependencies and generate a CSV file.
	$(GO) run github.com/google/go-licenses@latest report \
		--template '.github/license-3rdparty.tpl' \
		--ignore 'go.cipher.host/pkgdex' \
		'go.cipher.host/pkgdex' > LICENSE-3rdparty.csv

clean: # Cleans cache files from tests and deletes any build output.
	$(RM) -rf $(OUTDIR) ${ASSETSDIR}/css cover.out

.PHONY: all check pre-commit commit build build/linux build/css \
	build/doc install development tidy fmt lint lint/styles \
	lint/licenses fix/styles vulnerabilities test test/coverage \
	licenses clean
