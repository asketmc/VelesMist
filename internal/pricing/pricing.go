// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package pricing

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	apperrors "github.com/asketmc/VelesMist/internal/errors"
	"github.com/asketmc/VelesMist/internal/inventory"
)

type PriceInput struct {
	BuyerPriceCents int64  `json:"buyer_price_cents,omitempty"`
	LowestPrice     string `json:"lowest_price,omitempty"`
	MedianPrice     string `json:"median_price,omitempty"`
}

type Price struct {
	BuyerPriceCents int64  `json:"buyer_price_cents"`
	Source          string `json:"source"`
}

type PriceMap map[string]Price

type Options struct {
	ThresholdCents int64
	FeeBasisPoints int64
}

type PricedItem struct {
	AppID              int    `json:"appid"`
	Name               string `json:"name"`
	MarketHashName     string `json:"market_hash_name"`
	Count              int    `json:"count"`
	Tradable           bool   `json:"tradable"`
	MarketURL          string `json:"market_url"`
	PriceStatus        string `json:"price_status"`
	BuyerPriceCents    int64  `json:"buyer_price_cents,omitempty"`
	SellerReceiveCents int64  `json:"seller_receive_cents,omitempty"`
	TotalReceiveCents  int64  `json:"total_receive_cents,omitempty"`
	PriceSource        string `json:"price_source,omitempty"`
	Candidate          bool   `json:"candidate"`
}

type Analysis struct {
	Items      []PricedItem
	Candidates []PricedItem
}

func LoadPriceMap(raw map[string]PriceInput) (PriceMap, error) {
	prices := make(PriceMap, len(raw))
	for name, input := range raw {
		cents := input.BuyerPriceCents
		var err error
		if cents == 0 && input.LowestPrice != "" {
			cents, err = ParseMoneyToCents(input.LowestPrice)
		}
		if cents == 0 && input.MedianPrice != "" {
			cents, err = ParseMoneyToCents(input.MedianPrice)
		}
		if err != nil {
			return nil, apperrors.Wrap(apperrors.InvalidInput, "parse price for "+name, err)
		}
		if cents <= 0 {
			continue
		}
		prices[name] = Price{BuyerPriceCents: cents, Source: "cache"}
	}
	return prices, nil
}

func Analyze(items []inventory.AggregatedItem, prices PriceMap, opts Options) Analysis {
	if opts.FeeBasisPoints < 0 {
		opts.FeeBasisPoints = 0
	}
	result := Analysis{Items: make([]PricedItem, 0, len(items))}
	for _, item := range items {
		row := PricedItem{
			AppID:          item.AppID,
			Name:           item.Name,
			MarketHashName: item.MarketHashName,
			Count:          item.Count,
			Tradable:       item.Tradable,
			MarketURL:      MarketURL(item.AppID, item.MarketHashName),
			PriceStatus:    "missing",
		}
		if price, ok := prices[item.MarketHashName]; ok {
			seller := BuyerToSellerCents(price.BuyerPriceCents, opts.FeeBasisPoints)
			row.PriceStatus = "priced"
			row.BuyerPriceCents = price.BuyerPriceCents
			row.SellerReceiveCents = seller
			row.TotalReceiveCents = seller * int64(item.Count)
			row.PriceSource = price.Source
			row.Candidate = seller >= opts.ThresholdCents
		}
		result.Items = append(result.Items, row)
		if row.Candidate {
			result.Candidates = append(result.Candidates, row)
		}
	}
	sort.SliceStable(result.Candidates, func(i, j int) bool {
		if result.Candidates[i].TotalReceiveCents == result.Candidates[j].TotalReceiveCents {
			return result.Candidates[i].MarketHashName < result.Candidates[j].MarketHashName
		}
		return result.Candidates[i].TotalReceiveCents > result.Candidates[j].TotalReceiveCents
	})
	return result
}

func BuyerToSellerCents(buyerPriceCents int64, feeBasisPoints int64) int64 {
	if buyerPriceCents <= 0 {
		return 0
	}
	return buyerPriceCents * 10000 / (10000 + feeBasisPoints)
}

func ParseMoneyToCents(value string) (int64, error) {
	clean := strings.TrimSpace(value)
	clean = strings.TrimPrefix(clean, "$")
	clean = strings.TrimSuffix(clean, " USD")
	clean = strings.ReplaceAll(clean, ",", ".")
	if clean == "" {
		return 0, fmt.Errorf("empty money value")
	}
	parts := strings.Split(clean, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid money value %q", value)
	}
	dollars, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || dollars < 0 {
		return 0, fmt.Errorf("invalid money value %q", value)
	}
	var cents int64
	if len(parts) == 2 {
		frac := parts[1]
		if len(frac) == 1 {
			frac += "0"
		}
		if len(frac) > 2 {
			frac = frac[:2]
		}
		cents, err = strconv.ParseInt(frac, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid money value %q", value)
		}
	}
	return dollars*100 + cents, nil
}

func FormatCents(cents int64) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
}

func MarketURL(appID int, marketHashName string) string {
	return fmt.Sprintf("https://steamcommunity.com/market/listings/%d/%s", appID, url.PathEscape(marketHashName))
}
