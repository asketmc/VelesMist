// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package steam

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	apperrors "github.com/asketmc/VelesMist/internal/errors"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Options struct {
	BaseURL string
	Timeout time.Duration
}

func NewClient(opts Options) Client {
	if opts.BaseURL == "" {
		opts.BaseURL = "https://steamcommunity.com"
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 15 * time.Second
	}
	return Client{
		baseURL: strings.TrimRight(opts.BaseURL, "/"),
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
	}
}

func (c Client) FetchInventory(ctx context.Context, steamID string, appID int, contextID string) ([]byte, error) {
	endpoint, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.InvalidInput, "parse Steam base URL", err)
	}
	endpoint.Path = fmt.Sprintf("/inventory/%s/%d/%s", url.PathEscape(steamID), appID, url.PathEscape(contextID))
	q := endpoint.Query()
	q.Set("l", "english")
	q.Set("count", "5000")
	endpoint.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.InvalidInput, "build Steam request", err)
	}
	req.Header.Set("User-Agent", "VelesMist/0.1 (+https://github.com/asketmc/VelesMist)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, apperrors.Wrap(apperrors.NetworkTimeout, "Steam request timed out", err)
		}
		return nil, apperrors.Wrap(apperrors.Upstream, "Steam request failed", err)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if readErr != nil {
		return nil, apperrors.Wrap(apperrors.Upstream, "read Steam response", readErr)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, apperrors.New(apperrors.RateLimited, "Steam rate limited the request")
	}
	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusForbidden {
		return nil, apperrors.New(apperrors.Upstream, fmt.Sprintf("Steam inventory is private or unavailable (HTTP %d)", resp.StatusCode))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apperrors.New(apperrors.Upstream, fmt.Sprintf("Steam returned HTTP %d", resp.StatusCode))
	}
	return body, nil
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}
