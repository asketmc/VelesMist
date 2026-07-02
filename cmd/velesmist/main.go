// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	apperrors "github.com/asketmc/VelesMist/internal/errors"
	"github.com/asketmc/VelesMist/internal/cache"
	"github.com/asketmc/VelesMist/internal/config"
	"github.com/asketmc/VelesMist/internal/inventory"
	"github.com/asketmc/VelesMist/internal/pricing"
	"github.com/asketmc/VelesMist/internal/report"
	"github.com/asketmc/VelesMist/internal/steam"
	"github.com/asketmc/VelesMist/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: velesmist <scan|version> [options]")
		return apperrors.ExitInvalidInput
	}

	switch args[0] {
	case "version":
		return runVersion(stdout)
	case "scan":
		return runScan(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		return apperrors.ExitInvalidInput
	}
}

func runVersion(stdout io.Writer) int {
	info := version.Get()
	fmt.Fprintf(stdout, "velesmist %s commit=%s build_date=%s dirty=%s\n", info.Version, info.Commit, info.BuildDate, info.Dirty)
	return apperrors.ExitSuccess
}

func runScan(args []string, stdout io.Writer, stderr io.Writer) int {
	cfg, err := config.ParseScan(args)
	if err != nil {
		fmt.Fprintf(stderr, "invalid input: %v\n", err)
		return apperrors.ExitCode(err)
	}

	priceMap, err := loadPrices(cfg.PriceCache)
	if err != nil {
		fmt.Fprintf(stderr, "invalid price cache: %v\n", err)
		return apperrors.ExitCode(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	payload, err := fetchInventory(ctx, cfg)
	if err != nil {
		fmt.Fprintf(stderr, "scan failed: %v\n", err)
		return apperrors.ExitCode(err)
	}

	parsed, err := inventory.ParseSteamInventory(payload)
	if err != nil {
		fmt.Fprintf(stderr, "inventory parse failed: %v\n", err)
		return apperrors.ExitCode(err)
	}

	aggregated := inventory.AggregateMarketable(parsed)
	analysis := pricing.Analyze(aggregated, priceMap, pricing.Options{
		ThresholdCents: cfg.MinPriceCents,
		FeeBasisPoints: cfg.FeeBasisPoints,
	})

	result := report.BuildScanResult(report.ScanInput{
		SteamID:        cfg.SteamID,
		AppID:          cfg.AppID,
		ContextID:      cfg.ContextID,
		Currency:       cfg.Currency,
		ThresholdCents: cfg.MinPriceCents,
		GeneratedAt:    time.Now().UTC(),
		Items:          analysis.Items,
		Candidates:     analysis.Candidates,
	})

	switch cfg.Format {
	case config.FormatJSON:
		err = report.WriteJSON(stdout, result)
	case config.FormatTable:
		err = report.WriteTable(stdout, result)
	default:
		err = apperrors.New(apperrors.InvalidInput, "unsupported output format")
	}
	if err != nil {
		fmt.Fprintf(stderr, "report failed: %v\n", err)
		return apperrors.ExitCode(err)
	}

	return apperrors.ExitSuccess
}

func fetchInventory(ctx context.Context, cfg config.ScanConfig) ([]byte, error) {
	cacheKey := cache.InventoryKey(cfg.SteamID, cfg.AppID, cfg.ContextID)
	store := cache.NewStore(cfg.CacheFile)
	if !cfg.NoCache {
		if body, ok, err := store.GetValid(cacheKey, time.Now().UTC()); err == nil && ok {
			return body, nil
		}
	}

	client := steam.NewClient(steam.Options{
		BaseURL: cfg.SteamBaseURL,
		Timeout: cfg.Timeout,
	})
	payload, err := client.FetchInventory(ctx, cfg.SteamID, cfg.AppID, cfg.ContextID)
	if err != nil {
		return nil, err
	}

	if !cfg.NoCache {
		_ = store.Put(cacheKey, payload, time.Now().UTC(), cfg.CacheTTL)
	}
	return payload, nil
}

func loadPrices(path string) (pricing.PriceMap, error) {
	if path == "" {
		return pricing.PriceMap{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.InvalidInput, "open price cache", err)
	}
	defer file.Close()

	var raw map[string]pricing.PriceInput
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&raw); err != nil {
		return nil, apperrors.Wrap(apperrors.InvalidInput, "decode price cache", err)
	}
	return pricing.LoadPriceMap(raw)
}
