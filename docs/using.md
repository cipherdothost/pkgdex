<!--
SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>

SPDX-License-Identifier: CC-BY-SA-4.0
-->

# Using the service

This document explains how to use **pkgdex**, both as a user consuming
packages and as an administrator managing the service.

## For users

### Importing packages

To use a package "hosted" on **pkgdex** in your Go project, use `go get`
as usual, but replace the package repository URL with
`pkg.your-domain.com`:

```bash
go get pkg.your-domain.com/package-name@version
```

To import a package in your Go code, use the `import` statement as
usual:

```go
package main

import "pkg.your-domain.com/package-name"

func main() {
    // Do something with the package.
}
```

Of course, replace `pkg.your-domain.com` with the `baseURL` from your
configuration file.

### Search packages

The search feature allows you to find packages by:

- Name.
- Description.
- License.
- Keywords.

To search, visit `https://pkg.your-domain.com`, click the search bar at
the top, enter your search terms, and press Enter.

### RSS feed

You can stay updated with new packages and updates by subscribing to the
RSS feed at `https://pkg.your-domain.com/feed.xml`.

## For administrators

### Managing packages

#### Adding packages

1. Edit the configuration file:

```bash
sudo vim '/etc/pkgdex/config.json'
```

2. Add a package entry:
```json
{
  "packages": [
    {
      "name": "new-package",
      "description": "Package description",
      "version": "1.0.0",
      "branch": "trunk",
      "repository": "https://git.example.com/new-package",
      "license": "MIT"
    }
  ]
}
```

3. Restart the service:
```bash
sudo systemctl restart pkgdex
```

#### Generating usage examples

To generate formatted usage examples for packages, create a Go file
containing the example code and run the following command:

```bash
pkgdexctl generate-usage '/path/to/example.go'
```

The command should output the example code as a valid JSON string which
you can paste into the `usage` field of the package configuration.

#### Updating package versions

1. Update version in `config.json`:

```json
{
  "version": "1.1.0"
}
```

2. Restart the service:

```bash
sudo systemctl restart pkgdex
```

### Monitoring usage

#### View download statistics

**pkgdex** tries to keep track of how many times a package has been
downloaded using `go get`. The tracking is naive and can be inaccurate,
but it's better than nothing ¯\\\_(ツ)_/¯.

To view download statistics, run use cURL to access the API:

```bash
curl --silent \
  --header "Authorization: Bearer your-api-key" \
  'https://pkg.your-domain.com/meta/downloads'
```

#### Check service health

```bash
curl --silent \
  --header "Authorization: Bearer your-api-key" \
  'https://pkg.your-domain.com/meta/health'
```

### Common tasks

#### Service management

```bash
# Start the service.
sudo systemctl start pkgdex

# Stop the service.
sudo systemctl stop pkgdex

# Restart the service.
sudo systemctl restart pkgdex

# View service status.
sudo systemctl status pkgdex
```

#### Log Management

```bash
# View service logs.
sudo journalctl -u 'pkgdex'

# Follow logs in real-time.
sudo journalctl -u 'pkgdex' -f

# View logs since last boot.
sudo journalctl -u 'pkgdex' -b
```

#### Cache management

Assuming you're using NGINX as your reverse proxy and using our example
configuration, you can clear the NGINX cache as follows:

```bash
sudo rm -rf '/var/cache/nginx/pkgdex_cache/*'
sudo systemctl reload nginx
```

#### Backup management

We plan on adding a backup feature to the service in the future, but for
now, here's how you can do it manually:

```bash
# Stop the service.
sudo systemctl stop pkgdex

# Create a manual backup.
sudo tar czf "/var/backups/pkgdex-$(date +%Y%m%d).tar.gz" \
    '/etc/pkgdex' \
    '/var/lib/pkgdex'

# Start the service.
sudo systemctl start pkgdex
```

You'd probably want to automate this process with a `systemd` timer or a
cron job, but that's beyond the scope of this document.

To restore you'd do the following:

```bash
sudo systemctl stop pkgdex
sudo tar xzf '/var/backups/pkgdex-20250123.tar.gz' -C '/'
sudo systemctl start pkgdex
```
