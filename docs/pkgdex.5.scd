pkgdex(5)

# NAME

pkgdex - configuration file format for the *pkgdex* service

# DESCRIPTION

The *pkgdex* service uses a JSON configuration file to define its behavior. This
manual page describes the configuration file format and available options.

The default configuration file location is _/etc/pkgdex/config.json_.

# FILE FORMAT

The configuration file is a JSON document with four main sections:

- *service*: General service configuration.
- *server*: HTTP server settings.
- *database*: Database configuration.
- *packages*: List of packages with their metadata.

# SECTIONS

## SERVICE SECTION

The service section configures general service behavior:

```
{
  "service": {
    "seo": {
      "title": "Package Index",
      "description": "A package index for Go modules",
      "image": "https://example.com/image.png",
      "imageAlt": "Package index logo",
      "locale": "en_US",
      "publisher": "https://facebook.com/example",
      "twitter": "@example"
    },
    "name": "Example",
    "description": "A package index for Go modules",
    "contact": "hello@example.com",
    "homepage": "https://pkg.example.com",
    "baseURL": "pkg.example.com",
    "privacyPolicy": "https://pkg.example.com/privacy",
    "packagesPerPage": 10,
    "archive": false
  }
}
```

The following fields are available:

*seo*
	*title*
		Page title for SEO purposes. Optional.

	*description*
		Meta description for SEO purposes.

	*image*
		URL to the social media preview image.

	*imageAlt*
		Alt text for the social media preview image.

	*locale*
		Content locale in the IETF BCP 47 format.

	*publisher*
		Facebook page URL for the publisher.

	*twitter*
		Twitter handle for the publisher.

*name*
	Service name displayed in the web interface. Defaults to _pkgdex_.

*contact*
	Contact email address for the service administrator. Required.

*homepage*
	Service homepage URL. Required.

*baseURL*
	Base URL for package imports (e.g., pkg.example.com). Required.

*privacyPolicy*
	URL to the service's privacy policy. If no URL is given, a default privacy
	policy is used.

*archive*
	Enable archiving pages to the Wayback Machine. Optional. Defaults to _false_.

*packagesPerPage*
	Number of packages to display per page in listings. Optional. Defaults to _10_.

## SERVER SECTION

The server section configures the HTTP server:

```
{
  "server": {
    "address": "127.0.0.1:1997",
    "pid": "/var/run/pkgdex/pkgdex.pid",
    "readTimeout": "5s",
    "writeTimeout": "10s",
    "idleTimeout": "30s"
  }
}
```

The following fields are available:

*address*
	Server listen address and port. Optional. Defaults to _:1997_.

*pid*
	Path to PID file. Optional. Defaults to _/var/run/pkgdex.pid_.

*readTimeout*
	Maximum duration for reading the entire request. Optional. Defaults to _5s_.

*writeTimeout*
	Maximum duration before timing out writes of the response. Optional. Defaults
	to _10s_.

*idleTimeout*
	Maximum amount of time to wait for the next request. Optional. Defaults to
	_30s_.

## DATABASE SECTION

The database section configures the _bbolt_ database:

```
{
  "database": {
    "path": "/var/lib/pkgdex/database.db",
    "indexPath": "/var/lib/pkgdex/index"
  }
}
```

The following fields are available:

*path*
	Path to the database file. Optional. Defaults to
	_/var/lib/pkgdex/database.db_.

*indexPath*
	Path to the search index directory. Optional. Defaults to
	_/var/lib/pkgdex/index_.

## PACKAGES SECTION

The packages section contains a list of packages and their metadata:

```
{
  "packages": [
    {
      "name": "example",
      "description": "Example package",
      "version": "1.0.0",
      "branch": "trunk",
      "repository": "https://git.example.com/example.git",
      "license": "MIT",
      "usage": "// Optional example code",
      "hidden": false
    }
  ]
}
```

Each package entry supports the following fields:

*name*
	Package name. Required.

*description*
	Package description. Required.

*version*
	Package version in semver format (x.y.z). Required.

*branch*
	Git branch name. Required.

*repository*
	Git repository URL. Required. Must be HTTPS.

	*Supported hosts:*
		- github.com
		- gitlab.com
		- git.sr.ht

*license*
	Package license identifier. Required.

*usage*
	Example usage code. Optional. Must be generated using _pkgdexctl --generate-usage_.

*hidden*
	Whether to hide the package from listings.

# SECURITY

## TLS configuration

TLS certificates are managed using _systemd-creds_. Before starting the service,
encrypt your certificates.

## API keys

API keys protect sensitive routes and are managed using _systemd-creds_.

# FILES

*/etc/pkgdex/config.json*
	Default configuration file location.

*/etc/pkgdex/credentials/*
	Directory containing encrypted credentials.

*/var/lib/pkgdex/*
	Database and search index directory.

*/var/run/pkgdex/*
	Runtime files directory.

# EXAMPLES

Minimal configuration:

```
{
  "service": {
    "name": "My Package Index",
    "contact": "admin@example.com",
    "homepage": "https://pkg.example.com",
    "baseURL": "pkg.example.com",
    "privacyPolicy": "https://example.com/privacy"
  },
  "packages": [
    {
      "name": "example",
      "description": "Example package",
      "version": "1.0.0",
      "branch": "trunk",
      "repository": "https://github.com/example/example",
      "license": "MIT"
    }
  ]
}
```

Complete configuration with all options:

```
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
    "name": "My Package Index",
    "contact": "admin@example.com",
    "homepage": "https://pkg.example.com",
    "baseURL": "pkg.example.com",
    "privacyPolicy": "https://example.com/privacy",
    "archive": true,
    "packagesPerPage": 20
  },
  "server": {
    "address": "127.0.0.1:1997",
    "pid": "/var/run/pkgdex/pkgdex.pid",
    "readTimeout": "5s",
    "writeTimeout": "10s",
    "idleTimeout": "30s"
  },
  "database": {
    "path": "/var/lib/pkgdex/database.db",
    "indexPath": "/var/lib/pkgdex/index"
  },
  "packages": [
    {
      "name": "example",
      "description": "Example package",
      "version": "1.0.0",
      "branch": "trunk",
      "repository": "https://github.com/example/example",
      "license": "MIT",
      "usage": "// Generated with pkgdexctl --generate-usage",
      "hidden": false
    }
  ]
}
```

# SEE ALSO

*pkgdexctl*(1), *systemd-creds*(1)

# AUTHORS

Maintained by The Cipher Host Team <team@cipher.host>.
