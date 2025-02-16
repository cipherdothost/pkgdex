// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package sitemap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/version"
)

const (
	// xmlns is the XML namespace for sitemaps.
	xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

// Cache represents a cached sitemap.
type Cache struct {
	sitemap []byte
}

// NewCache generates and caches a sitemap for the given configuration.
func NewCache(cfg *config.Config, versionMgr *version.Manager, lastModified time.Time) (*Cache, error) {
	baseURL := strings.TrimSuffix(cfg.Service.Homepage, "/")

	urlset := URLSet{
		XMLNS: xmlns,
		URLs: []URL{
			{
				Location:     baseURL,
				LastModified: lastModified.UTC().Format(time.DateOnly),
			},
		},
	}

	if cfg.Service.PrivacyPolicy == "" {
		urlset.URLs = append(urlset.URLs, URL{
			Location:     baseURL + "/privacy",
			LastModified: lastModified.UTC().Format(time.DateOnly),
		})
	}

	for _, pkg := range cfg.Packages {
		if pkg.Hidden {
			continue
		}

		timestamp := versionMgr.GetLastUpdate(pkg, time.Now())

		urlset.URLs = append(urlset.URLs, URL{
			Location:     baseURL + "/" + pkg.Name,
			LastModified: timestamp.Format(time.DateOnly),
		})
	}

	var (
		estimatedSize = len(xml.Header) + (len(cfg.Packages) * 500)
		buf           = bytes.NewBuffer(make([]byte, 0, estimatedSize))
	)

	buf.WriteString(xml.Header)

	encoder := xml.NewEncoder(buf)
	encoder.Indent("", "  ")

	if err := encoder.Encode(urlset); err != nil {
		return nil, fmt.Errorf("failed to generate sitemap: %w", err)
	}

	return &Cache{
		sitemap: buf.Bytes(),
	}, nil
}

// Get returns the cached sitemap content.
func (c *Cache) Get() []byte {
	return c.sitemap
}
