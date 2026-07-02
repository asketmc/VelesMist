// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package inventory

import (
	"os"
	"testing"
)

func TestParseSteamInventoryMatchesAssetsToDescriptions(t *testing.T) {
	payload := mustReadFixture(t)
	items, err := ParseSteamInventory(payload)
	if err != nil {
		t.Fatalf("ParseSteamInventory returned error: %v", err)
	}
	if len(items) != 4 {
		t.Fatalf("len(items) = %d, want 4", len(items))
	}
	if items[0].MarketHashName != "Golden Moonfall" || !items[0].Marketable {
		t.Fatalf("unexpected first item: %+v", items[0])
	}
}

func TestAggregateMarketableExcludesUnmarketableAndCountsDuplicates(t *testing.T) {
	items, err := ParseSteamInventory(mustReadFixture(t))
	if err != nil {
		t.Fatalf("ParseSteamInventory returned error: %v", err)
	}
	aggregated := AggregateMarketable(items)
	byName := map[string]AggregatedItem{}
	for _, item := range aggregated {
		byName[item.MarketHashName] = item
	}
	if len(byName) != 2 {
		t.Fatalf("marketable item count = %d, want 2", len(byName))
	}
	if byName["Golden Moonfall"].Count != 2 {
		t.Fatalf("Golden Moonfall count = %d, want 2", byName["Golden Moonfall"].Count)
	}
	if _, ok := byName["Unmarketable Relic"]; ok {
		t.Fatal("unmarketable item was aggregated")
	}
}

func FuzzParseSteamInventory(f *testing.F) {
	body, err := os.ReadFile("testdata/dota_inventory.json")
	if err != nil {
		f.Fatalf("read fixture: %v", err)
	}
	f.Add(string(body))
	f.Add(`{"success":1,"assets":[],"descriptions":[]}`)
	f.Fuzz(func(t *testing.T, payload string) {
		_, _ = ParseSteamInventory([]byte(payload))
	})
}

func mustReadFixture(t testing.TB) []byte {
	t.Helper()
	body, err := os.ReadFile("testdata/dota_inventory.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return body
}
