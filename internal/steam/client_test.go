// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package steam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	apperrors "github.com/asketmc/VelesMist/internal/errors"
)

func TestFetchInventoryUsesMockHTTPServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/inventory/76561198000000000/570/2" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("count") != "5000" {
			t.Fatalf("count query = %s", r.URL.Query().Get("count"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":1,"assets":[],"descriptions":[]}`))
	}))
	defer server.Close()

	client := NewClient(Options{BaseURL: server.URL, Timeout: time.Second})
	body, err := client.FetchInventory(context.Background(), "76561198000000000", 570, "2")
	if err != nil {
		t.Fatalf("FetchInventory error: %v", err)
	}
	if string(body) != `{"success":1,"assets":[],"descriptions":[]}` {
		t.Fatalf("body = %s", body)
	}
}

func TestFetchInventoryMapsRateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := NewClient(Options{BaseURL: server.URL, Timeout: time.Second})
	_, err := client.FetchInventory(context.Background(), "76561198000000000", 570, "2")
	if got := apperrors.KindOf(err); got != apperrors.RateLimited {
		t.Fatalf("KindOf(err) = %s, want %s", got, apperrors.RateLimited)
	}
}

func TestFetchInventoryMapsPrivateOrUnavailableInventory(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			}))
			defer server.Close()

			client := NewClient(Options{BaseURL: server.URL, Timeout: time.Second})
			_, err := client.FetchInventory(context.Background(), "76561197987179126", 570, "2")
			if got := apperrors.KindOf(err); got != apperrors.Upstream {
				t.Fatalf("KindOf(err) = %s, want %s", got, apperrors.Upstream)
			}
			if err == nil || !strings.Contains(err.Error(), "private or unavailable") {
				t.Fatalf("error = %v, want private/unavailable message", err)
			}
		})
	}
}
