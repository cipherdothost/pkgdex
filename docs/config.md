<!--
SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>

SPDX-License-Identifier: CC-BY-SA-4.0
-->

# Configuration reference

## File location

The default configuration file location is `/etc/pkgdex/config.json`.
The service expects a JSON configuration file with the following main
sections:

- `service`: General service configuration.
- `server`: HTTP server settings.
- `database`: Database configuration.
- `packages`: List of packages with their metadata.

## Service configuration

The service section configures general service behavior:

```json
{
  "service": {
    "seo": {
      "title": "Package Index",
      "description": "A package index for Go modules",
      "image": "https://example.com/image.png",
      "imageAlt": "Package Index Logo",
      "locale": "en_US",
      "publisher": "https://facebook.com/example",
      "twitter": "@example"
    },
    "name": "Example",
    "description": "A package index for Go modules",
    "contact": "hello@example.com",
    "homepage": "https://example.com",
    "baseURL": "pkg.example.com",
    "privacyPolicy": "https://example.com/privacy",
    "packagesPerPage": 10,
    "archive": false
  }
}
```

| Field | Description | Default | Required |
|-------|-------------|---------|----------|
| `name` | Service name | `pkgdex` | No |
| `description` | Service description | Package index description | No |
| `contact` | Contact email address | - | Yes |
| `homepage` | Service homepage URL | - | Yes |
| `baseURL` | Base URL for package imports | - | Yes |
| `privacyPolicy` | Privacy policy URL | - | No |
| `packagesPerPage` | Number of packages per page | `10` | No |
| `archive` | Enable Wayback Machine archiving | `false` | No |

### SEO configuration

| Field | Description | Required |
|-------|-------------|----------|
| `title` | Homepage title | No |
| `description` | Meta description | No |
| `image` | Social media image URL | No |
| `imageAlt` | Image alt text | No |
| `locale` | Site locale | No |
| `publisher` | Facebook page URL | No |
| `twitter` | Twitter handle | No |

### Server configuration

The server section configures the HTTP server:

```json
{
  "server": {
    "address": "pkg.example.com:1997",
    "pid": "/var/run/pkgdex/pkgdex.pid",
    "readTimeout": "5s",
    "writeTimeout": "10s",
    "idleTimeout": "30s"
  }
}
```

| Field | Description | Default | Required |
|-------|-------------|---------|----------|
| `address` | Server listen address | `:1997` | No |
| `pid` | Path to PID file | `/var/run/pkgdex/pkgdex.pid` | No |
| `readTimeout` | Read timeout duration | `5s` | No |
| `writeTimeout` | Write timeout duration | `10s` | No |
| `idleTimeout` | Idle timeout duration | `30s` | No |

### Database configuration

The database section configures the bbolt database:

```json
{
  "database": {
    "path": "/var/lib/pkgdex/pkgdex.db",
    "indexPath": "/var/lib/pkgdex/index"
  }
}
```

| Field | Description | Default | Required |
|-------|-------------|---------|----------|
| `path` | Database file path | `/var/lib/pkgdex/pkgdex.db` | No |
| `indexPath` | Search index directory path | `/var/lib/pkgdex/index` | No |

### Package configuration

The packages section contains a list of packages and their metadata:

```json
{
  "packages": [
    {
      "name": "example",
      "description": "Example package",
      "version": "1.0.0",
      "branch": "trunk",
      "repository": "https://git.example.com/example.git",
      "license": "MIT",
      "usage": "// Optional example code generated with pkgdex --generate-usage",
      "hidden": false
    }
  ]
}
```

| Field | Description | Required | Format |
|-------|-------------|----------|--------|
| `name` | Package name | Yes | `^[a-zA-Z0-9][a-zA-Z0-9-_/.]*[a-zA-Z0-9]$` |
| `description` | Package description | Yes | - |
| `version` | Package version | Yes | `x.y.z` |
| `branch` | Git branch name | Yes | - |
| `repository` | Git repository URL | Yes | Valid URL |
| `license` | Package license | Yes | - |
| `usage` | Example usage code | No | - |
| `hidden` | Hide package from listings | No | `false` |

## Security

### TLS configuration

TLS certificate and private key are managed using `systemd` and `systemd-creds`. You must encrypt these credentials before starting the service:

```bash
# Create a secure temporary directory.
TMPDIR="$(mktemp -d -p /dev/shm)"
chmod 700 "${TMPDIR}"

# Copy or generate your TLS certificate and key to the temporary directory.
cp '/path/to/your/cert.pem' "${TMPDIR}/cert.pem"
cp '/path/to/your/key.pem' "${TMPDIR}/key.pem"

# Encrypt the credentials using systemd-creds.
sudo systemd-creds encrypt "${TMPDIR}/cert.pem" '/etc/pkgdex/cert.pem.creds'
sudo systemd-creds encrypt "${TMPDIR}/key.pem" '/etc/pkgdex/key.pem.creds'

# Clean up temporary files.
rm -rf "${TMPDIR}"
```

In your `systemd` service file, load the credentials using:

```ini
[Service]
LoadCredentialEncrypted=pkgdex-tlscertificate:/etc/pkgdex/cert.pem.creds
LoadCredentialEncrypted=pkgdex-tlskey:/etc/pkgdex/key.pem.creds
```

### API keys

The service uses an API key for protecting sensitive routes. The key is automatically generated if not set and must:
- Be at least 64 characters long.
- Start with `pkgdex_`.

To generate and securely store a custom API key:

```bash
# Create a secure temporary directory.
TMPDIR="$(mktemp -d -p /dev/shm)"
chmod 700 "${TMPDIR}"

# Generate the API key and store it in the temporary directory.
pkgdex generate-key > "${TMPDIR}/key"

# Encrypt the API key using systemd-creds.
sudo systemd-creds encrypt "${TMPDIR}/api.key" /etc/pkgdex/key.creds

# Clean up temporary files.
rm -rf "${TMPDIR}"
```

In your `systemd` service file, load the API key using:

```ini
[Service]
LoadCredentialEncrypted=pkgdex-key:/etc/pkgdex/key.creds
```

## Example configuration

A complete example configuration can be found at
[`../config/config.example.json`](../config/config.example.json).

You can also find example configuration files for:

- [`apparmor`](../config/apparmor/)
- [`nginx`](../config/nginx/)
- [`systemd`](../config/systemd/)

The configuration we use for production can be found in [the
go.cipher.host
repository](https://github.com/cipherdothost/go.cipher.host).
