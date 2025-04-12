<!--
SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>

SPDX-License-Identifier: CC0-1.0
-->

# Pkgdex

![A screenshot of the Pkgdex homepage showing two Go packages and their
metadata.](.github/assets/pkgdex-homepage-screenshot.png "Pkgdex
homepage")

**Pkgdex** is a powerful, yet simple, self-hosted service that provides
both a searchable index for your Go packages and [vanity import
path](https://pkg.go.dev/cmd/go#hdr-Remote_import_paths) handling. 

## Features

- Full-text search across packages.
- Vanity import path support.
- Download statistics.
- RSS feed for updates.
- XML sitemap generation.
- Wayback Machine integration.

A live version of the index is available at
[https://go.cipher.host/](https://go.cipher.host/).

## Installation

### From source

First install the dependencies:

- Go 1.24 or above.
- make.
- npm.
- [scdoc](https://git.sr.ht/~sircmpwn/scdoc).

Clone the repository, switch to the latest stable tag, then compile, and
install:

```bash
git clone 'https://github.com/cipherdothost/pkgdex.git'
cd 'pkgdex'
git checkout 'v1.0.0'
npm install
make
sudo make install
```

## Usage

- [Installation guide](docs/hosting.md)
- [Usage instructions](docs/using.md)
- [Configuration reference](docs/config.md)

## Contributing

Anyone can help make **pkgdex** better. Check out [the contribution
guidelines](CONTRIBUTING.md) and [the development
instructions](docs/development.md) for more information.

---

The work in this repository complies with [the REUSE
specification](https://reuse.software/spec-3.3/). While [the default
license is EUPL-1.2](LICENSE.md), individual files may be licensed
differently.

Please see the individual files for details and [the LICENSES
directory](LICENSES/) for a full list of licenses used in this
repository.
