// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/asketmc/VelesMist/internal/pricing"
)

const SchemaVersion = "velesmist.scan.v1"

type ScanInput struct {
	SteamID        string
	AppID          int
	ContextID      string
	Currency       string
	ThresholdCents int64
	GeneratedAt    time.Time
	Items          []pricing.PricedItem
	Candidates     []pricing.PricedItem
}

type ScanResult struct {
	SchemaVersion  string               `json:"schema_version"`
	GeneratedAt    time.Time            `json:"generated_at"`
	SteamID        string               `json:"steam_id"`
	AppID          int                  `json:"appid"`
	ContextID      string               `json:"contextid"`
	Currency       string               `json:"currency"`
	ThresholdCents int64                `json:"threshold_cents"`
	Items          []pricing.PricedItem `json:"items"`
	Candidates     []pricing.PricedItem `json:"candidates"`
	Summary        Summary              `json:"summary"`
}

type Summary struct {
	MarketableItems            int   `json:"marketable_items"`
	PricedItems                int   `json:"priced_items"`
	MissingPriceItems          int   `json:"missing_price_items"`
	SkippedItems               int   `json:"skipped_items"`
	CandidateItems             int   `json:"candidate_items"`
	EstimatedTotalGrossCents   int64 `json:"estimated_total_gross_cents"`
	EstimatedTotalFeeCents     int64 `json:"estimated_total_fee_cents"`
	EstimatedTotalReceiveCents int64 `json:"estimated_total_receive_cents"`
}

func BuildScanResult(input ScanInput) ScanResult {
	summary := Summary{
		MarketableItems: len(input.Items),
		CandidateItems:  len(input.Candidates),
	}
	for _, item := range input.Items {
		if item.PriceStatus == pricing.PriceStatusPriced {
			summary.PricedItems++
		}
		switch item.Recommendation {
		case pricing.RecommendationMissingPrice:
			summary.MissingPriceItems++
		case pricing.RecommendationSkip:
			summary.SkippedItems++
		}
	}
	for _, item := range input.Candidates {
		summary.EstimatedTotalGrossCents += item.TotalBuyerPriceCents
		summary.EstimatedTotalFeeCents += item.TotalEstimatedFeeCents
		summary.EstimatedTotalReceiveCents += item.TotalReceiveCents
	}
	return ScanResult{
		SchemaVersion:  SchemaVersion,
		GeneratedAt:    input.GeneratedAt,
		SteamID:        input.SteamID,
		AppID:          input.AppID,
		ContextID:      input.ContextID,
		Currency:       input.Currency,
		ThresholdCents: input.ThresholdCents,
		Items:          input.Items,
		Candidates:     input.Candidates,
		Summary:        summary,
	}
}

func WriteJSON(w io.Writer, result ScanResult) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func WriteTable(w io.Writer, result ScanResult) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ITEM\tCOUNT\tGROSS\tFEE\tYOU_RECEIVE\tTOTAL\tCONFIDENCE\tRECOMMENDATION\tREASONS\tMARKET_URL")
	for _, item := range result.Items {
		gross := "-"
		fee := "-"
		seller := "-"
		total := "-"
		if item.PriceStatus == pricing.PriceStatusPriced {
			gross = result.Currency + " " + pricing.FormatCents(item.BuyerPriceCents)
			fee = result.Currency + " " + pricing.FormatCents(item.EstimatedFeeCents)
			seller = result.Currency + " " + pricing.FormatCents(item.SellerReceiveCents)
			total = result.Currency + " " + pricing.FormatCents(item.TotalReceiveCents)
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			item.MarketHashName,
			item.Count,
			gross,
			fee,
			seller,
			total,
			item.Confidence,
			item.Recommendation,
			strings.Join(item.ReasonCodes, ","),
			item.MarketURL,
		)
	}
	fmt.Fprintf(tw, "\nSell recommendations: %d Estimated gross: %s %s Estimated fees: %s %s Estimated receive: %s %s\n",
		result.Summary.CandidateItems,
		result.Currency,
		pricing.FormatCents(result.Summary.EstimatedTotalGrossCents),
		result.Currency,
		pricing.FormatCents(result.Summary.EstimatedTotalFeeCents),
		result.Currency,
		pricing.FormatCents(result.Summary.EstimatedTotalReceiveCents),
	)
	return tw.Flush()
}
