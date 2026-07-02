// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package pricing

import (
	"testing"

	"github.com/asketmc/VelesMist/internal/inventory"
)

func TestParseMoneyToCents(t *testing.T) {
	tests := map[string]int64{
		"$12.34":  1234,
		"12,30":   1230,
		"5":       500,
		"0.99 USD": 99,
	}
	for input, want := range tests {
		got, err := ParseMoneyToCents(input)
		if err != nil {
			t.Fatalf("ParseMoneyToCents(%q) error: %v", input, err)
		}
		if got != want {
			t.Fatalf("ParseMoneyToCents(%q) = %d, want %d", input, got, want)
		}
	}
}

func TestAnalyzeFiltersBySellerReceiveThreshold(t *testing.T) {
	items := []inventory.AggregatedItem{
		{AppID: 570, Name: "Golden Moonfall", MarketHashName: "Golden Moonfall", Count: 2, Tradable: true},
		{AppID: 570, Name: "Cheap", MarketHashName: "Cheap", Count: 1, Tradable: true},
	}
	prices := PriceMap{
		"Golden Moonfall": {BuyerPriceCents: 1234, Source: "cache"},
		"Cheap":          {BuyerPriceCents: 400, Source: "cache"},
	}
	result := Analyze(items, prices, Options{ThresholdCents: 500, FeeBasisPoints: 1500})
	if len(result.Candidates) != 1 {
		t.Fatalf("candidate count = %d, want 1", len(result.Candidates))
	}
	got := result.Candidates[0]
	if got.MarketHashName != "Golden Moonfall" {
		t.Fatalf("candidate = %s, want Golden Moonfall", got.MarketHashName)
	}
	if got.SellerReceiveCents != 1073 || got.TotalReceiveCents != 2146 {
		t.Fatalf("seller=%d total=%d, want seller=1073 total=2146", got.SellerReceiveCents, got.TotalReceiveCents)
	}
}

func TestLoadPriceMapUsesLowestThenMedian(t *testing.T) {
	got, err := LoadPriceMap(map[string]PriceInput{
		"A": {LowestPrice: "$2.50", MedianPrice: "$2.30"},
		"B": {MedianPrice: "$1.25"},
	})
	if err != nil {
		t.Fatalf("LoadPriceMap error: %v", err)
	}
	if got["A"].BuyerPriceCents != 250 || got["B"].BuyerPriceCents != 125 {
		t.Fatalf("unexpected prices: %+v", got)
	}
}
