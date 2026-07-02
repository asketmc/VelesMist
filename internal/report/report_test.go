// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package report

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/asketmc/VelesMist/internal/pricing"
)

func TestWriteJSONGolden(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteJSON(&buf, testResult()); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	assertGolden(t, "testdata/scan.json.golden", buf.String())
}

func TestWriteTableGolden(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteTable(&buf, testResult()); err != nil {
		t.Fatalf("WriteTable error: %v", err)
	}
	assertGolden(t, "testdata/scan.table.golden", buf.String())
}

func testResult() ScanResult {
	items := []pricing.PricedItem{
		{
			AppID:                  570,
			Name:                   "Golden Moonfall",
			MarketHashName:         "Golden Moonfall",
			Count:                  2,
			Tradable:               true,
			MarketURL:              "https://steamcommunity.com/market/listings/570/Golden%20Moonfall",
			PriceStatus:            pricing.PriceStatusPriced,
			BuyerPriceCents:        1234,
			EstimatedFeeCents:      161,
			SellerReceiveCents:     1073,
			TotalBuyerPriceCents:   2468,
			TotalEstimatedFeeCents: 322,
			TotalReceiveCents:      2146,
			PriceSource:            "cache",
			Recommendation:         pricing.RecommendationSell,
			Candidate:              true,
		},
		{
			AppID:          570,
			Name:           "Jagged Honor | Blade",
			MarketHashName: "Jagged Honor | Blade",
			Count:          1,
			Tradable:       true,
			MarketURL:      "https://steamcommunity.com/market/listings/570/Jagged%20Honor%20%7C%20Blade",
			PriceStatus:    pricing.PriceStatusMissing,
			Recommendation: pricing.RecommendationMissingPrice,
		},
	}
	return BuildScanResult(ScanInput{
		SteamID:        "76561198000000000",
		AppID:          570,
		ContextID:      "2",
		Currency:       "USD",
		ThresholdCents: 500,
		GeneratedAt:    time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC),
		Items:          items,
		Candidates:     items[:1],
	})
}

func assertGolden(t *testing.T, path string, got string) {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if got != string(body) {
		t.Fatalf("golden mismatch\n--- got ---\n%s\n--- want ---\n%s", got, body)
	}
}
