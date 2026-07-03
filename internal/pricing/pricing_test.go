// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package pricing

import (
	"strings"
	"testing"

	"github.com/asketmc/VelesMist/internal/inventory"
)

func TestParseMoneyToCents(t *testing.T) {
	tests := map[string]int64{
		"$12.34":   1234,
		"12,30":    1230,
		"5":        500,
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
		"Cheap":           {BuyerPriceCents: 400, Source: "cache"},
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
	if got.EstimatedFeeCents != 161 || got.TotalEstimatedFeeCents != 322 {
		t.Fatalf("fee=%d total_fee=%d, want fee=161 total_fee=322", got.EstimatedFeeCents, got.TotalEstimatedFeeCents)
	}
	if got.Recommendation != RecommendationSell {
		t.Fatalf("recommendation = %s, want %s", got.Recommendation, RecommendationSell)
	}
	if got.Confidence != ConfidenceMedium {
		t.Fatalf("confidence = %s, want %s", got.Confidence, ConfidenceMedium)
	}
	if len(got.ReasonCodes) == 0 {
		t.Fatal("expected reason codes")
	}
	if result.Items[1].Recommendation != RecommendationSkip {
		t.Fatalf("cheap recommendation = %s, want %s", result.Items[1].Recommendation, RecommendationSkip)
	}
	if result.Items[1].ReasonCodes[len(result.Items[1].ReasonCodes)-1] != ReasonBelowMinNet {
		t.Fatalf("cheap reason codes = %v, want last %s", result.Items[1].ReasonCodes, ReasonBelowMinNet)
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

func TestLoadPriceCacheV1(t *testing.T) {
	input := `{
	  "schema_version": "velesmist.price-cache.v1",
	  "currency": "USD",
	  "prices": {
	    "Golden Moonfall": {
	      "buyer_price_cents": 1234,
	      "source": "manual"
	    }
	  }
	}`
	prices, meta, err := LoadPriceCache(strings.NewReader(input))
	if err != nil {
		t.Fatalf("LoadPriceCache error: %v", err)
	}
	if meta.SchemaVersion != PriceCacheSchemaVersion || meta.Currency != "USD" || meta.Count != 1 {
		t.Fatalf("unexpected metadata: %+v", meta)
	}
	if prices["Golden Moonfall"].BuyerPriceCents != 1234 || prices["Golden Moonfall"].Source != "manual" {
		t.Fatalf("unexpected price: %+v", prices["Golden Moonfall"])
	}
}

func TestLoadPriceCacheRejectsWrongSchema(t *testing.T) {
	input := `{"schema_version":"v0","currency":"USD","prices":{}}`
	if _, _, err := LoadPriceCache(strings.NewReader(input)); err == nil {
		t.Fatal("expected schema error")
	}
}
