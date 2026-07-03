// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/asketmc/VelesMist/internal/cache"
	"github.com/asketmc/VelesMist/internal/errors"
	"github.com/asketmc/VelesMist/internal/pricing"
	"github.com/asketmc/VelesMist/internal/report"
)

func TestRunUsageUnknownAndVersion(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantExit   int
		wantStdout string
		wantStderr string
	}{
		{
			name:       "usage",
			args:       nil,
			wantExit:   errors.ExitInvalidInput,
			wantStderr: "usage: velesmist <scan|prices|ui|version> [options]",
		},
		{
			name:       "unknown",
			args:       []string{"bogus"},
			wantExit:   errors.ExitInvalidInput,
			wantStderr: "unknown command: bogus",
		},
		{
			name:       "version",
			args:       []string{"version"},
			wantExit:   errors.ExitSuccess,
			wantStdout: "velesmist",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := run(tt.args, &stdout, &stderr)
			if code != tt.wantExit {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
			}
			if tt.wantStdout != "" && !strings.Contains(stdout.String(), tt.wantStdout) {
				t.Fatalf("stdout=%q, want %q", stdout.String(), tt.wantStdout)
			}
			if tt.wantStderr != "" && !strings.Contains(stderr.String(), tt.wantStderr) {
				t.Fatalf("stderr=%q, want %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestRunScanRejectsInvalidInput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"scan", "--format", "yaml"}, &stdout, &stderr)

	if code != errors.ExitInvalidInput {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "invalid input") {
		t.Fatalf("stderr=%q, want invalid input", stderr.String())
	}
}

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

func TestPricesCommandValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "missing subcommand", args: []string{"prices"}, want: "usage: velesmist prices <template>"},
		{name: "unknown subcommand", args: []string{"prices", "bogus"}, want: "unknown prices command: bogus"},
		{name: "template positional arg", args: []string{"prices", "template", "extra"}, want: "invalid input"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := run(tt.args, &stdout, &stderr)
			if code != errors.ExitInvalidInput {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("stderr=%q, want %q", stderr.String(), tt.want)
			}
		})
	}
}

func TestPricesTemplateWritesFileAndHonorsForce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prices.json")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run([]string{"prices", "template", "--output", path}, &stdout, &stderr)
	if code != errors.ExitSuccess {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout=%q, want no stdout for file output", stdout.String())
	}
	if _, _, err := pricing.LoadPriceCache(mustOpen(t, path)); err != nil {
		t.Fatalf("written price cache did not parse: %v", err)
	}

	stderr.Reset()
	code = run([]string{"prices", "template", "--output", path}, &stdout, &stderr)
	if code != errors.ExitInvalidInput {
		t.Fatalf("exit=%d stderr=%s, want existing file failure", code, stderr.String())
	}

	stderr.Reset()
	code = run([]string{"prices", "template", "--output", path, "--force"}, &stdout, &stderr)
	if code != errors.ExitSuccess {
		t.Fatalf("force exit=%d stderr=%s", code, stderr.String())
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

func TestWriteScanResultFormatsAndRejectsUnknownFormat(t *testing.T) {
	result := report.ScanResult{
		SchemaVersion: report.SchemaVersion,
		Currency:      "USD",
		Items: []pricing.PricedItem{{
			MarketHashName:     "Golden Moonfall",
			Count:              1,
			PriceStatus:        pricing.PriceStatusPriced,
			BuyerPriceCents:    1234,
			EstimatedFeeCents:  161,
			SellerReceiveCents: 1073,
			TotalReceiveCents:  1073,
			Recommendation:     pricing.RecommendationSell,
			Confidence:         pricing.ConfidenceMedium,
			ReasonCodes:        []string{pricing.ReasonMarketable},
			MarketURL:          "https://steamcommunity.com/market/listings/570/Golden%20Moonfall",
		}},
		Summary: report.Summary{CandidateItems: 1},
	}

	var jsonOut bytes.Buffer
	if err := writeScanResult(&jsonOut, "json", result); err != nil {
		t.Fatalf("write json: %v", err)
	}
	if !json.Valid(jsonOut.Bytes()) {
		t.Fatalf("invalid json output: %s", jsonOut.String())
	}

	var tableOut bytes.Buffer
	if err := writeScanResult(&tableOut, "table", result); err != nil {
		t.Fatalf("write table: %v", err)
	}
	if !strings.Contains(tableOut.String(), "Golden Moonfall") || !strings.Contains(tableOut.String(), "Sell recommendations") {
		t.Fatalf("unexpected table output: %s", tableOut.String())
	}

	if err := writeScanResult(&bytes.Buffer{}, "xml", result); err == nil {
		t.Fatal("expected unknown format error")
	}
}

func TestRunUIReportsListenFailure(t *testing.T) {
	listener, err := netListenLocalhost()
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runUI([]string{"--addr", listener.Addr().String(), "--open=false"}, &stdout, &stderr)
	if code != errors.ExitInvalidInput {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "listen on UI address") {
		t.Fatalf("stderr=%q, want listen failure", stderr.String())
	}
}

func mustOpen(t *testing.T, path string) *os.File {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	t.Cleanup(func() { _ = file.Close() })
	return file
}

func netListenLocalhost() (net.Listener, error) {
	return net.Listen("tcp", "127.0.0.1:0")
}
