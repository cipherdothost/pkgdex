// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package wayback

import (
	"context"

	"go.cipher.host/pkgdex/internal/config"
)

// Archive manages the preservation of package URLs to the Wayback Machine.
type Archive struct {
	client *Client
}

// NewArchive returns a new instance of Archive.
func NewArchive(client *Client) *Archive {
	if client == nil {
		client = NewClient(nil)
	}

	return &Archive{
		client: client,
	}
}

// Save saves the given URI to the Wayback Machine if it is not already archived.
func (a *Archive) Save(ctx context.Context, uri string) error {
	archived, err := a.client.IsArchived(ctx, uri)
	if err != nil {
		return err
	}

	if archived {
		return nil
	}

	return a.client.Archive(ctx, uri)
}

// SavePackages saves the given packages to the Wayback Machine if they are not
// already archived.
func (a *Archive) SavePackages(ctx context.Context, cfg *config.Config) error {
	for _, pkg := range cfg.Packages {
		uri := PackageURI("https://"+cfg.Service.BaseURL, pkg.Name)

		if err := a.Save(ctx, uri); err != nil {
			return err
		}
	}

	return nil
}
