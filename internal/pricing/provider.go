// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package pricing

import (
	"context"
	"os"

	"github.com/asketmc/VelesMist/internal/domain"
	apperrors "github.com/asketmc/VelesMist/internal/errors"
	"github.com/asketmc/VelesMist/internal/inventory"
)

type FilePriceProvider struct{}

func (FilePriceProvider) FetchPrices(_ context.Context, req domain.PriceRequest) (PriceMap, error) {
	if req.Path == "" {
		return PriceMap{}, nil
	}
	file, err := os.Open(req.Path)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.InvalidInput, "open price cache", err)
	}
	defer file.Close()
	prices, _, err := LoadPriceCache(file)
	return prices, err
}

type Scorer struct{}

func (Scorer) Score(_ context.Context, items []inventory.AggregatedItem, prices PriceMap, opts Options) (Analysis, error) {
	return Analyze(items, prices, opts), nil
}
