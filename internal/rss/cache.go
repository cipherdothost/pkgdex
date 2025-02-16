// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package rss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/version"
)

// Cache represents a cached RSS feed.
type Cache struct {
	feed []byte
}

// NewCache generates and caches an RSS feed for the given configuration.
func NewCache(cfg *config.Config, versionMgr *version.Manager) (*Cache, error) {
	baseURL := strings.TrimSuffix(cfg.Service.Homepage, "/")

	feed := &Feed{
		Version: "2.0",
		Channel: Channel{
			Title:       cfg.Service.Name,
			Link:        baseURL,
			Description: "Latest Go packages from " + cfg.Service.Name + ".",
			Language:    "en-us",
			Generator:   cfg.Service.Name,
			Items:       make([]Item, 0, len(cfg.Packages)),
		},
	}

	var lastBuildDate time.Time

	for _, pkg := range cfg.Packages {
		if pkg.Hidden {
			continue
		}

		timestamp := versionMgr.GetLastUpdate(pkg, time.Now())

		if timestamp.After(lastBuildDate) {
			lastBuildDate = timestamp
		}

		item := Item{
			Title:       pkg.Name,
			Link:        baseURL + "/" + pkg.Name,
			Description: pkg.Description,
			PubDate:     timestamp.Format(time.RFC1123Z),
			GUID:        baseURL + "/" + pkg.Name,
		}

		feed.Channel.Items = append(feed.Channel.Items, item)
	}

	feed.Channel.LastBuildDate = lastBuildDate.Format(time.RFC1123Z)

	buf := new(bytes.Buffer)
	buf.WriteString(xml.Header)

	encoder := xml.NewEncoder(buf)
	encoder.Indent("", "  ")

	if err := encoder.Encode(feed); err != nil {
		return nil, fmt.Errorf("failed to generate RSS feed: %w", err)
	}

	return &Cache{
		feed: buf.Bytes(),
	}, nil
}

// Get returns the cached feed content.
func (c *Cache) Get() []byte {
	return c.feed
}
