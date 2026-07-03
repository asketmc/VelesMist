// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"time"

	"github.com/asketmc/VelesMist/internal/domain"
	"github.com/asketmc/VelesMist/internal/inventory"
	"github.com/asketmc/VelesMist/internal/pricing"
	"github.com/asketmc/VelesMist/internal/report"
)

type InventoryProvider interface {
	FetchInventory(ctx context.Context, req domain.InventoryRequest) (inventory.Inventory, error)
}

type PriceProvider interface {
	FetchPrices(ctx context.Context, req domain.PriceRequest) (pricing.PriceMap, error)
}

type Scorer interface {
	Score(ctx context.Context, items []inventory.AggregatedItem, prices pricing.PriceMap, opts pricing.Options) (pricing.Analysis, error)
}

type Scanner struct {
	InventoryProvider InventoryProvider
	PriceProvider     PriceProvider
	Scorer            Scorer
	Clock             func() time.Time
}

func (s Scanner) Scan(ctx context.Context, req domain.ScanRequest) (report.ScanResult, error) {
	inv, err := s.InventoryProvider.FetchInventory(ctx, req.Inventory)
	if err != nil {
		return report.ScanResult{}, err
	}
	prices, err := s.PriceProvider.FetchPrices(ctx, req.Prices)
	if err != nil {
		return report.ScanResult{}, err
	}
	analysis, err := s.Scorer.Score(ctx, inventory.AggregateMarketable(inv.Items), prices, pricing.Options{
		ThresholdCents: req.ThresholdCents,
		FeeBasisPoints: req.FeeBasisPoints,
	})
	if err != nil {
		return report.ScanResult{}, err
	}
	now := time.Now().UTC()
	if s.Clock != nil {
		now = s.Clock().UTC()
	}
	return report.BuildScanResult(report.ScanInput{
		SteamID:        req.Inventory.SteamID,
		AppID:          req.Inventory.AppID,
		ContextID:      req.Inventory.ContextID,
		Currency:       req.Currency,
		ThresholdCents: req.ThresholdCents,
		GeneratedAt:    now,
		Items:          analysis.Items,
		Candidates:     analysis.Candidates,
	}), nil
}
