// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/asketmc/VelesMist/internal/cache"
	"github.com/asketmc/VelesMist/internal/errors"
	"github.com/asketmc/VelesMist/internal/pricing"
	"github.com/asketmc/VelesMist/internal/report"
)

func TestPricesTemplateWritesPriceCacheV1(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"prices", "template"}, &stdout, &stderr)
	if code != errors.ExitSuccess {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}

	prices, meta, err := pricing.LoadPriceCache(bytes.NewReader(stdout.Bytes()))
	if err != nil {
		t.Fatalf("template did not parse as price cache: %v", err)
	}
	if meta.SchemaVersion != pricing.PriceCacheSchemaVersion || meta.Count == 0 || len(prices) == 0 {
		t.Fatalf("unexpected template metadata=%+v prices=%+v", meta, prices)
	}
}

func TestScanUsesLocalCacheAndPriceCacheJSON(t *testing.T) {
	dir := t.TempDir()
	inventoryCachePath := filepath.Join(dir, "inventory-cache.json")
	priceCachePath := filepath.Join(dir, "price-cache.json")

	body, err := os.ReadFile(filepath.Join("..", "..", "internal", "inventory", "testdata", "dota_inventory.json"))
	if err != nil {
		t.Fatalf("read inventory fixture: %v", err)
	}
	store := cache.NewStore(inventoryCachePath)
	key := cache.InventoryKey("76561198000000000", 570, "2")
	if err := store.Put(key, body, time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC), 24*time.Hour); err != nil {
		t.Fatalf("write cache fixture: %v", err)
	}

	priceCache := `{
	  "schema_version": "velesmist.price-cache.v1",
	  "currency": "USD",
	  "prices": {
	    "Golden Moonfall": {"buyer_price_cents": 1234, "source": "test"},
	    "Jagged Honor | Blade": {"buyer_price_cents": 450, "source": "test"}
	  }
	}`
	if err := os.WriteFile(priceCachePath, []byte(priceCache), 0o600); err != nil {
		t.Fatalf("write price cache fixture: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run([]string{
		"scan",
		"--steam-id", "76561198000000000",
		"--format", "json",
		"--cache-file", inventoryCachePath,
		"--price-cache", priceCachePath,
		"--min-price", "5.00",
	}, &stdout, &stderr)
	if code != errors.ExitSuccess {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}

	var got report.ScanResult
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("decode scan JSON: %v\n%s", err, stdout.String())
	}
	if got.Summary.CandidateItems != 1 || got.Summary.SkippedItems != 1 || got.Summary.MissingPriceItems != 0 {
		t.Fatalf("unexpected summary: %+v", got.Summary)
	}
	if got.Items[0].Recommendation != pricing.RecommendationSell || got.Items[1].Recommendation != pricing.RecommendationSkip {
		t.Fatalf("unexpected recommendations: %+v", got.Items)
	}
	if got.Items[0].MarketURL == "" || got.Items[0].EstimatedFeeCents == 0 {
		t.Fatalf("expected market URL and fee in first item: %+v", got.Items[0])
	}
}

func TestScanUsesFixtureWithoutNetworkOrSteamID(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run([]string{
		"scan",
		"--fixture", filepath.Join("..", "..", "internal", "inventory", "testdata", "dota_inventory.json"),
		"--format", "json",
	}, &stdout, &stderr)
	if code != errors.ExitSuccess {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var got report.ScanResult
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("decode scan JSON: %v\n%s", err, stdout.String())
	}
	if got.AppID != 570 || got.ContextID != "2" {
		t.Fatalf("unexpected fixture report target: appid=%d contextid=%s", got.AppID, got.ContextID)
	}
	if got.Summary.MarketableItems != 2 || got.Summary.MissingPriceItems != 2 {
		t.Fatalf("unexpected fixture summary: %+v", got.Summary)
	}
}
