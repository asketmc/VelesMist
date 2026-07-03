// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package steam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/asketmc/VelesMist/internal/cache"
	"github.com/asketmc/VelesMist/internal/domain"
)

func TestInventoryProviderUsesCacheBeforeHTTP(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	body := []byte(`{"success":1,"assets":[],"descriptions":[]}`)
	key := cache.InventoryKey("76561198000000000", 570, "2")
	if err := cache.NewStore(path).Put(key, body, time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC), time.Hour); err != nil {
		t.Fatalf("write cache: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("HTTP server should not be called on cache hit")
	}))
	defer server.Close()

	provider := NewInventoryProvider(InventoryProviderOptions{
		Client:    NewClient(Options{BaseURL: server.URL, Timeout: time.Second}),
		CacheFile: path,
		CacheTTL:  time.Hour,
		Clock:     func() time.Time { return time.Date(2026, 7, 3, 0, 1, 0, 0, time.UTC) },
	})
	inv, err := provider.FetchInventory(context.Background(), domain.InventoryRequest{
		SteamID:   "76561198000000000",
		AppID:     570,
		ContextID: "2",
	})
	if err != nil {
		t.Fatalf("FetchInventory error: %v", err)
	}
	if len(inv.Items) != 0 {
		t.Fatalf("items = %d, want 0", len(inv.Items))
	}
}
