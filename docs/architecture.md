<!--
SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>

SPDX-License-Identifier: CC-BY-SA-4.0
-->

# Architecture

This document describes the high-level architecture of **pkgdex**. If
you want to familiarize yourself with the code base, you are just in the
right place!

## Bird's eye view

At its core, **pkgdex** is a Go service that provides two main
functionalities:

1. A searchable index of Go packages with metadata, download statistics,
   and documentation links.
2. A vanity import path handler that enables custom import paths for Go
   packages.

The service is designed to be self-hosted and operates as a standalone
HTTP server. It maintains all package metadata in a local
[etcd-io/bbolt](https://github.com/etcd-io/bbolt) database and provides
a web interface for browsing and searching packages. 

## Code map

This section describes the main components of **pkgdex** and how they
interact with each other.

**TODO:** Write this section.
