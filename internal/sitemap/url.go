// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package sitemap

import "encoding/xml"

// URLSet represents the root element of a sitemap.
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL represents a URL entry in the sitemap.
type URL struct {
	Location     string `xml:"loc"`
	LastModified string `xml:"lastmod"`
}
