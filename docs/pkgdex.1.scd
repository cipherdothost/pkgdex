pkgdexctl(1)

# NAME

pkgdexctl - control the HTTP server for the *pkgdex* service

# SYNOPSIS

*pkgdexctl* [OPTIONS]

# DESCRIPTION

*pkgdexctl* is a control application for the *pkgdex* service, which provides a
package index for Go packages.

The control application can be used to start and stop the server, generate API
keys, show package download statistics, and generate usage examples for package
documentation.

# OPTIONS

*-h*, *--help*
	Show help message and quit.

*-v*, *--version*
	Show version number and quit.

# COMMANDS

*start* [OPTIONS]
	Starts the server using the configuration specified with *--config*.

	Options are:

	*-c*, *--config* <FILE>
		Path to the configuration file. Defaults to _/etc/pkgdex/config.json_.

*stop* [OPTIONS]
	Stops a running server instance.

	Options are:

	*-c*, *--config* <FILE>
		Path to the configuration file. Defaults to _/etc/pkgdex/config.json_.

*generate-usage* <FILE>
	Convert a Go source file to a usage example string for the configuration
	file. The file must have a _.go_ extension.

*generate-key*
	Generate a valid API key for protecting sensitive routes.

# CONFIGURATION

The configuration file is a JSON file with the following main sections:

- *service*: General service configuration.
- *server*: HTTP server settings.
- *database*: Database configuration.
- *packages*: List of packages with their metadata.

See *pkgdex*(5) for detailed configuration file documentation.

# EXIT STATUS

*0*
	Success.

*1*
	Unknown error.

*2*
	Output error.

*3*
	Usage error.

*4*
	Interrupted.

*5*
	Configuration error.

*6*
	Database error.

*7*
	Server error.

*8*
	Filesystem error.

*9*
	Network error.

*10*
	Authentication error.

*11*
	Validation error.

# EXAMPLES

*Example 1. Start the server:*
	$ pkgdexctl start --config '/etc/pkgdex/config.json'

*Example 2. Stop the server:*
	$ pkgdexctl stop --config '/etc/pkgdex/config.json'

*Example 3. Generate an API key:*
	$ pkgdexctl generate-key

*Example 4. Generate a usage example string from a Go file:*
	$ pkgdexctl generate-usage '/path/to/example.go'

# SECURITY

Credentials are managed using systemd-creds and must be set up before starting
the service. See *pkgdex*(5) for details.

# SEE ALSO

*pkgdex*(5), *systemd-creds*(1)

# AUTHORS

Maintained by The Cipher Host Team <team@cipher.host>.

# BUGS

Submit bug reports online at: <https://github.com/cipherdothost/pkgdex>
