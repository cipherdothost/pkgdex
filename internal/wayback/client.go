// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package wayback

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
	"git.sr.ht/~jamesponddotco/xstd-go/xstrings"
	jsoniter "github.com/json-iterator/go"
	"go.cipher.host/cmdkit"
	"golang.org/x/time/rate"
)

// ErrRequest is returned if an HTTP request fails for whatever reason.
const ErrRequest cmdkit.Error = "Wayback Machine request failed; the service may be unavailable or rate-limiting"

const (
	archiveBaseURL string = "https://web.archive.org"
	checkBaseURL   string = "http://archive.org/wayback/available"
)

// Doer represents an HTTP client that can execute requests.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type (
	// Client represents a client to the Wayback Machine.
	Client struct {
		http    Doer
		limiter *rate.Limiter
	}

	// Response represents a response from the Wayback Machine API.
	Response struct {
		response *http.Response
		body     []byte
	}
)

// New returns a new instance of Client.
func NewClient(httpClient Doer) *Client {
	if httpClient == nil {
		httpClient = xhttp.NewClient(15 * time.Second)
	}

	return &Client{
		http:    httpClient,
		limiter: rate.NewLimiter(rate.Limit(0.5), 1),
	}
}

// Archive attempts to save the given URI to the Wayback Machine.
func (c *Client) Archive(ctx context.Context, uri string) error {
	reqURI := xstrings.JoinWithSeparator("/", archiveBaseURL, "save", uri)

	if _, err := url.ParseRequestURI(reqURI); err != nil {
		return fmt.Errorf("%w", err)
	}

	form := url.Values{}
	form.Set("url", uri)
	form.Set("capture_all", "on")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURI, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.do(req)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// IsArchived checks if the given URL is already archived in the Wayback Machine.
func (c *Client) IsArchived(ctx context.Context, uri string) (bool, error) {
	reqURI := fmt.Sprintf("%s?url=%s", checkBaseURL, url.QueryEscape(uri))

	if _, err := url.ParseRequestURI(reqURI); err != nil {
		return false, fmt.Errorf("%w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURI, http.NoBody)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	var result struct {
		ArchivedSnapshots struct { //nolint:revive // we'll write a proper client at a later date which won't have this issue
			Closest struct { //nolint:revive // see above`
				Available bool `json:"available"`
			} `json:"closest"`
		} `json:"archived_snapshots"`
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	if err = json.NewDecoder(bytes.NewReader(resp.body)).Decode(&result); err != nil {
		return false, fmt.Errorf("parsing response: %w", err)
	}

	return result.ArchivedSnapshots.Closest.Available, nil
}

// Do executes an HTTP request and returns the response.
func (c *Client) do(req *http.Request) (*Response, error) {
	if err := c.limiter.Wait(req.Context()); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	defer func() {
		if err = xhttp.DrainResponseBody(resp); err != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode > http.StatusNoContent {
		return nil, fmt.Errorf("%w: with status code: %d", ErrRequest, resp.StatusCode)
	}

	var buffer *bytes.Buffer

	if resp.ContentLength > 0 {
		buffer = bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	} else {
		buffer = bytes.NewBuffer(make([]byte, 0, 1024))
	}

	_, err = io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Response{
		response: resp,
		body:     buffer.Bytes(),
	}, nil
}
