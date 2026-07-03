// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"testing"
	"time"

	"github.com/asketmc/VelesMist/internal/domain"
	"github.com/asketmc/VelesMist/internal/inventory"
	"github.com/asketmc/VelesMist/internal/pricing"
)

type fakeInventoryProvider struct {
	inventory inventory.Inventory
	err       error
}

func (p fakeInventoryProvider) FetchInventory(context.Context, domain.InventoryRequest) (inventory.Inventory, error) {
	return p.inventory, p.err
}

type fakePriceProvider struct {
	prices pricing.PriceMap
	err    error
}

func (p fakePriceProvider) FetchPrices(context.Context, domain.PriceRequest) (pricing.PriceMap, error) {
	return p.prices, p.err
}

func TestScannerBuildsReportThroughProviders(t *testing.T) {
	scanner := Scanner{
		InventoryProvider: fakeInventoryProvider{inventory: inventory.Inventory{Items: []inventory.Item{
			{
				AppID:          570,
				Name:           "Golden Moonfall",
				MarketHashName: "Golden Moonfall",
				Amount:         2,
				Marketable:     true,
				Tradable:       true,
			},
		}}},
		PriceProvider: fakePriceProvider{prices: pricing.PriceMap{
			"Golden Moonfall": {BuyerPriceCents: 1234, Source: "test"},
		}},
		Scorer: pricing.Scorer{},
		Clock:  func() time.Time { return time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC) },
	}
	got, err := scanner.Scan(context.Background(), domain.ScanRequest{
		Inventory: domain.InventoryRequest{
			SteamID:   "76561198000000000",
			Game:      "dota2",
			AppID:     570,
			ContextID: "2",
		},
		Currency:       "USD",
		ThresholdCents: 500,
		FeeBasisPoints: 1500,
	})
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if got.Summary.CandidateItems != 1 || got.Summary.EstimatedTotalReceiveCents != 2146 {
		t.Fatalf("unexpected summary: %+v", got.Summary)
	}
	if got.Items[0].Recommendation != pricing.RecommendationSell {
		t.Fatalf("recommendation = %s, want %s", got.Items[0].Recommendation, pricing.RecommendationSell)
	}
}
