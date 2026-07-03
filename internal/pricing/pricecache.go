// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package pricing

import (
	"encoding/json"
	"io"

	apperrors "github.com/asketmc/VelesMist/internal/errors"
)

const PriceCacheSchemaVersion = "velesmist.price-cache.v1"

type PriceCacheFile struct {
	SchemaVersion string                `json:"schema_version"`
	Currency      string                `json:"currency"`
	Prices        map[string]PriceInput `json:"prices"`
}

type PriceCacheMetadata struct {
	SchemaVersion string
	Currency      string
	Count         int
}

func LoadPriceCache(r io.Reader) (PriceMap, PriceCacheMetadata, error) {
	var file PriceCacheFile
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&file); err != nil {
		return nil, PriceCacheMetadata{}, apperrors.Wrap(apperrors.InvalidInput, "decode price cache", err)
	}
	if file.SchemaVersion != PriceCacheSchemaVersion {
		return nil, PriceCacheMetadata{}, apperrors.New(apperrors.InvalidInput, "unsupported price cache schema version")
	}
	if file.Prices == nil {
		return nil, PriceCacheMetadata{}, apperrors.New(apperrors.InvalidInput, "price cache prices object is required")
	}
	prices, err := LoadPriceMap(file.Prices)
	if err != nil {
		return nil, PriceCacheMetadata{}, err
	}
	return prices, PriceCacheMetadata{
		SchemaVersion: file.SchemaVersion,
		Currency:      file.Currency,
		Count:         len(prices),
	}, nil
}

func WritePriceCacheTemplate(w io.Writer) error {
	template := PriceCacheFile{
		SchemaVersion: PriceCacheSchemaVersion,
		Currency:      "USD",
		Prices: map[string]PriceInput{
			"Example Market Hash Name": {
				BuyerPriceCents: 1234,
				Source:          "manual",
			},
			"Example Steam Lowest Price": {
				LowestPrice: "$1.23",
				Source:      "manual-steam-market-check",
			},
		},
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(template)
}
