<!--
SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>

SPDX-License-Identifier: CC-BY-SA-4.0
-->

# Development

This document provides technical guidance for developers who want to
contribute code to **pkgdex**. It complements the
[../CONTRIBUTING.md](../CONTRIBUTING.md) file by focusing specifically
on the technical aspects of development.

Before you begin, please review the
[../CONTRIBUTING.md](../CONTRIBUTING.md) file to understand the overall
contribution guidelines and our code of conduct.

## Table of Contents

1. [Development environment](#development-environment-setup)
2. [Building and Testing](#building-and-testing)
3. [Contribution Workflow](#contribution-workflow)
4. [Coding Standards](#coding-standards)
5. [Documentation](#documentation)
6. [License Compliance](#license-compliance)

## Development Environment

### Prerequisites

- Go 1.24 or above
- Node/NPM
- make
- scdoc

### Setting up your environment

1. Fork the repository on GitHub.

2. Clone your fork locally:
```bash
git clone 'https://github.com/YOUR-USERNAME/pkgdex.git'
cd 'pkgdex'
```

3. Create a development configuration file:
```bash
mkdir -p '.dev'
cp 'config/dev-config.example.json' '.dev/config.json'
```

4. Edit the `.dev/config.json` file to set appropriate development
   values. You'll want to edit the following fields:

   - `server.pid`
   - `database.path`
   - `database.indexPath`

You may also want to edit the following fields:

   - `service.homepage`
   - `service.baseURL`
   - `server.address`

### Creating a self-signed certificate for development

Since **pkgdex** requires TLS, you'll need to generate a self-signed
certificate for local development:

```bash
openssl req -x509 -newkey rsa:4096 -keyout '.dev/pkgdex-tlskey' -out '.dev/pkgdex-tlscertificate' -days 365 -nodes -subj '/CN=localhost'
```

## Building and Testing

### Building the project

To build the project locally, run:

```bash
npm install
make
```

This will compile the binary to `build/pkgdexctl`.

For development, it's best to use the `development` target with your
development configuration in `.dev`:

```bash
make development
```

By default, this will start the service on `https://localhost:8080`, but
you can change this in the `.dev/config.json` file by modifying the
`server.address` field.  You'll need to accept the self-signed
certificate warning.

> [!NOTE]
> If you modify the server address, you may also want to modify the
> `service.homepage` and `service.baseURL` fields.

### Running tests

To run the test suite:

```bash
make test
```

For test output with coverage information:

```bash
make test/coverage
```

## Contribution Workflow

As said before, review the [../CONTRIBUTING.md](../CONTRIBUTING.md) file
to understand the overall contribution guidelines and our code of
conduct.

### Creating a feature or fix

1. Ensure there's an issue for your work:
   - Check [existing
     issues](https://github.com/cipherdothost/pkgdex/issues) first.
   - If none exists, create one clearly describing the bug or feature.

2. Create a new branch for your work in your fork:
```bash
git checkout -b 'feature/your-feature-name'
# or
git checkout -b 'fix/issue-description'
```

3. Make your changes, following the coding standards.

4. Write tests for your changes.

5. Update documentation if necessary.

6. Update the `Unreleased` section of the
   [../CHANGELOG.md](../CHANGELOG.md) file.

### Submitting pull requests

1. Push your branch to your fork:
```bash
git push origin 'feature/your-feature-name'
```

2. Create a pull request against the `trunk` branch of the main
   repository.

3. In your PR description:
   - Clearly describe the changes.
   - Reference the issue it resolves.
   - Note any breaking changes or migration steps.

4. Wait for code review feedback and address any comments.

## Coding Standards

### Go code style

We like to follow both the Uber and Google Go coding standards, with a
preference for the Uber one when they conflict with each other, and we
expect your contribution to do the same.

- [Uber style guide](https://github.com/uber-go/guide/blob/master/style.md).
- [Google style guide](https://google.github.io/styleguide/go/)

### Dependencies

As mentioned in the [../CONTRIBUTING.md](../CONTRIBUTING.md) file, avoid
adding new dependencies unless absolutely necessary.

If you need to add a new dependency:

1. Ensure it's well-maintained and widely used.
2. Check its license for compatibility.
3. Update the [../LICENSE-3rdparty.csv](../LICENSE-3rdparty.csv) file.
4. Run `make tidy` to update the `go.mod` and `go.sum` files.

## Documentation

### Code documentation

- Add godoc comments to all types, functions, and methods, including
  non-exported ones.
- Comments should explain the why behind the code, not just the what.
- Put the most important information first in the comment. Aim to make
  the first sentence a clear, concise summary.
  Focus on the intent and purpose of each part of the code.
- Use complete sentences, proper capitalization and punctuation in
  comments.

### Project documentation

Project documentation is stored in the `docs/` directory:

- `hosting.md`: Instructions for deploying the service.
- `config.md`: Configuration reference.
- `using.md`: Usage instructions.
- `development.md`: Development instructions.

When updating features or APIs, make corresponding updates to the
documentation. If adding new configuration options, update `config.md`
accordingly.

## License Compliance

Remember that all code must comply with [the REUSE
specification](https://reuse.software/spec-3.3/). The default license is
[EUPL-1.2](../LICENSE.md).

Ensure each new file needs the appropriate license header. You can use
the `reuse` tool to help with this:

```bash
pipx install reuse
make lint/licenses
```

---

Thank you for contributing to **pkgdex**! Your efforts help make this
project better for everyone.
