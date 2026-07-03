// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package steam

import (
	"context"
	"encoding/json"
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
		if r.URL.Query().Get("count") != inventoryPageSize {
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

func TestFetchInventoryPaginatesWithSafePageSize(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Query().Get("count") != inventoryPageSize {
			t.Fatalf("count query = %s", r.URL.Query().Get("count"))
		}
		switch requests {
		case 1:
			if got := r.URL.Query().Get("start_assetid"); got != "" {
				t.Fatalf("first start_assetid = %s", got)
			}
			_, _ = w.Write([]byte(`{
				"success":1,
				"assets":[{"appid":570,"contextid":"2","assetid":"1","classid":"10","instanceid":"0","amount":"1"}],
				"descriptions":[{"appid":570,"classid":"10","instanceid":"0","name":"Page One","market_hash_name":"Page One","marketable":1,"tradable":1}],
				"more_items":1,
				"last_assetid":"1",
				"total_inventory_count":2
			}`))
		case 2:
			if got := r.URL.Query().Get("start_assetid"); got != "1" {
				t.Fatalf("second start_assetid = %s, want 1", got)
			}
			_, _ = w.Write([]byte(`{
				"success":1,
				"assets":[{"appid":570,"contextid":"2","assetid":"2","classid":"20","instanceid":"0","amount":"1"}],
				"descriptions":[{"appid":570,"classid":"20","instanceid":"0","name":"Page Two","market_hash_name":"Page Two","marketable":1,"tradable":1}],
				"more_items":0,
				"total_inventory_count":2
			}`))
		default:
			t.Fatalf("unexpected request %d", requests)
		}
	}))
	defer server.Close()

	client := NewClient(Options{BaseURL: server.URL, Timeout: time.Second})
	body, err := client.FetchInventory(context.Background(), "76561198000000000", 570, "2")
	if err != nil {
		t.Fatalf("FetchInventory error: %v", err)
	}
	var got struct {
		Assets       []json.RawMessage `json:"assets"`
		Descriptions []json.RawMessage `json:"descriptions"`
		TotalCount   int               `json:"total_inventory_count"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode combined body: %v\n%s", err, body)
	}
	if requests != 2 || len(got.Assets) != 2 || len(got.Descriptions) != 2 || got.TotalCount != 2 {
		t.Fatalf("requests=%d combined=%+v body=%s", requests, got, body)
	}
}

func TestFetchInventoryRejectsNonAdvancingPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"success":1,
			"assets":[],
			"descriptions":[],
			"more_items":true,
			"last_assetid":"same"
		}`))
	}))
	defer server.Close()

	client := NewClient(Options{BaseURL: server.URL, Timeout: time.Second})
	_, err := client.FetchInventory(context.Background(), "76561198000000000", 570, "2")
	if got := apperrors.KindOf(err); got != apperrors.Upstream {
		t.Fatalf("KindOf(err) = %s, want %s", got, apperrors.Upstream)
	}
	if err == nil || !strings.Contains(err.Error(), "pagination did not advance") {
		t.Fatalf("error = %v, want pagination advance error", err)
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
			want := "private or unavailable"
			if status == http.StatusBadRequest {
				want = "rejected the inventory request"
			}
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("error = %v, want %q", err, want)
			}
		})
	}
}
