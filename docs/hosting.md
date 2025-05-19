<!--
SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>

SPDX-License-Identifier: CC-BY-SA-4.0
-->

# Hosting the service

This document provides instructions for hosting the **pkgdex** service
on a production server.

## System requirements

- Linux server with `systemd`. This guide was written with [openSUSE
  Tumbleweed](https://www.opensuse.org/) in mind, but any Linux
  distribution with `systemd` should work.
- NGINX, Caddy, or another reverse proxy.
- Go 1.24 or later, `make`, and `scdoc`.
- Root access or sudo privileges.

While not required, `apparmor` is recommended for extra security.

## Installation

**1. Create the system user and directories:**

```bash
# Create the system user.
useradd -r -s '/sbin/nologin' -U 'pkgdex'

# Create the required directories.
sudo mkdir -p '/etc/pkgdex/credentials'
sudo mkdir -p '/var/lib/pkgdex'
sudo mkdir -p '/var/run/pkgdex'
sudo mkdir -p '/var/cache/nginx/pkgdex_cache'

# Set proper ownership and permissions.
sudo chown -R 'pkgdex:pkgdex' '/etc/pkgdex'
sudo chown -R 'pkgdex:pkgdex' '/var/lib/pkgdex'
sudo chown -R 'pkgdex:pkgdex' '/var/run/pkgdex'
sudo chown -R 'root:root' '/var/cache/nginx/pkgdex_cache'
sudo chmod 750 '/etc/pkgdex'
sudo chmod 750 '/var/lib/pkgdex'
sudo chmod 750 '/var/run/pkgdex'
```

**2. Build and install the service:**

```bash
# Clone the repository.
git clone 'https://github.com/cipherdothost/pkgdex.git'
cd 'pkgdex'

# Switch to latest stable version.
git checkout 'v1.0.1'

# Build the service.
make

# Install the service.
sudo make install
```

**3. Configure the TLS certificate:**

The service requires a valid TLS certificate to work, so you'll need to
obtain and manage SSL/TLS certificates before proceeding with this step.

> **Important**: This documentation doesn't cover certificate
> acquisition and management. You should set up an automated
> certificate management solution such as:
>
> - [`lego`](https://go-acme.github.io/lego/)
> - [`certbot`](https://certbot.eff.org)
> - [`acme.sh`](https://github.com/acmesh-official/acme.sh)

Once you have your certificates:

```bash
# Generate a temporary directory for handling secrets.
TMPDIR="$(mktemp -d -p /dev/shm)"
chmod 700 "${TMPDIR}"

# Copy your TLS certificates.
cp '/path/to/your/cert.pem' "${TMPDIR}/cert.pem"
cp '/path/to/your/key.pem' "${TMPDIR}/key.pem"

# Encrypt the certificates using systemd-creds.
sudo systemd-creds encrypt "${TMPDIR}/cert.pem" '/etc/pkgdex/credentials/pkgdex-tlscertificate'
sudo systemd-creds encrypt "${TMPDIR}/key.pem" '/etc/pkgdex/credentials/pkgdex-tlskey'
```

**4. Generate and configure an API key:**

This is optional, but if you want to access the `/meta/` endpoints, you
will need an API key. If you don't care about that, skip this step as
the service will generate a random every time it starts.

```bash
# Generate the API key.
sudo pkgdexctl generate-key > "${TMPDIR}/key"

# Encrypt the API key.
sudo systemd-creds encrypt "${TMPDIR}/key" '/etc/pkgdex/credentials/pkgdex-key'

# Clean up.
rm -rf "${TMPDIR}"
```

**5. Configure the service:**

Create a `config.json` file in `/etc/pkgdex` with your configuration.
You can find an example configuration file at
[`../config/config.example.json`](../config/config.example.json) and a
configuration reference at [`config.md`](config.md).

Edit the file and replace the placeholders with your values.

**6. Install the `apparmor` profile:**

This step is optional, but implementing `apparmor` significantly
enhances your server's security through mandatory access control. 

- It restricts the service to only necessary system resources.
- It prevents unauthorized file system access.
- It helps mitigate damage from potential security breaches.
- It provides an additional layer of defense-in-depth.

If your system has `apparmor` available, and openSUSE Tumbleweed does,
follow these steps:

```bash
# First, verify AppArmor is available and running.
sudo aa-status

# Copy AppArmor profile.
sudo cp '/usr/local/share/pkgdex/apparmor/usr.local.bin.pkgdexctl' '/etc/apparmor.d/'

# Load the profile.
sudo apparmor_parser -r '/etc/apparmor.d/usr.local.bin.pkgdexctl'

# Verify the profile is loaded.
sudo aa-status | grep 'pkgdex'
```

If you choose not to use `apparmor`, you should implement alternative
security measures such as:
- Strict filesystem permissions.
- SELinux policies (if using a RedHat-based system).
- Comprehensive system auditing.
- Network access controls.

**7. Configure NGINX:**

While the service can be accessed over HTTPS directly, it is recommended
to use a reverse proxy such as NGINX or Caddy. Features such as caching,
rate limiting, and other security features are easier to implement with
a reverse proxy.

To set up NGINX as a reverse proxy, follow these steps:

```bash
# Copy the NGINX configuration example.
sudo cp '/usr/local/share/pkgdex/nginx/pkg.example.com' '/etc/nginx/vhosts.d/pkg.your-domain.com'

# Edit the configuration file.
sudo vim '/etc/nginx/vhosts.d/pkg.your-domain.com'

# Test and reload NGINX.
sudo nginx -t
sudo systemctl reload nginx
```

**8. Install and start the `systemd` service:**

```bash
# Copy the systemd service file.
sudo cp '/usr/local/share/pkgdex/systemd/pkgdex.service' '/etc/systemd/system/'

# Reload systemd.
sudo systemctl daemon-reload

# Enable and start the service.
sudo systemctl enable pkgdex
sudo systemctl start pkgdex
```

## Verification

**1. Check service status:**

```bash
sudo systemctl status pkgdex
```

**2. Check logs:**

```bash
sudo journalctl -u pkgdex
```

**3. Test the API:**

```bash
curl -sH 'Authorization: Bearer your-api-key' 'https://pkg.your-domain.com/meta/health'
```

## Security considerations

1. Ensure all directories have proper permissions (750 or more
   restrictive).
2. Regularly update TLS certificates.
3. Monitor logs for unauthorized access attempts.
4. Keep the system, the server, and Go updated.
5. Consider implementing additional network security measures (firewall
   rules, fail2ban, etc.).

## FAQ

**1. What should I backup?**

Everything, really, but the most important files would be:

   - **Database**: `/var/lib/pkgdex/database.db`
   - **Search index**: `/var/lib/pkgdex/index`
   - **Configuration**: `/etc/pkgdex/config.json`
   - **Credentials**: `/etc/pkgdex/credentials`

**2. A new version is out. How do I upgrade the service?**

```bash
# Clone the repository.
git clone 'https://github.com/cipherdothost/pkgdex.git'
cd 'pkgdex'

# Switch to latest stable version.
git checkout 'v{{VERSION}}'

# Build the service.
make

# Install the service again.
sudo make install

# Restart the service.
sudo systemctl restart pkgdex
```

**3. Where are my logs? How do I set up a log rotation policy?**

The service uses `systemd`'s journal for logging. Configure `journald`
according to your retention needs.

**4. The service won't start. What can I do?**

Most issues can be found in the logs. Try:

```bash
journalctl -u pkgdex -n 50
```
