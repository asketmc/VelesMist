// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/asketmc/VelesMist/internal/app"
	"github.com/asketmc/VelesMist/internal/config"
	"github.com/asketmc/VelesMist/internal/domain"
	apperrors "github.com/asketmc/VelesMist/internal/errors"
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
		fmt.Fprintln(stderr, "usage: velesmist <scan|prices|ui|version> [options]")
		return apperrors.ExitInvalidInput
	}

	switch args[0] {
	case "version":
		return runVersion(stdout)
	case "scan":
		return runScan(args[1:], stdout, stderr)
	case "prices":
		return runPrices(args[1:], stdout, stderr)
	case "ui":
		return runUI(args[1:], stdout, stderr)
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

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	result, err := scanWithConfig(ctx, cfg)
	if err != nil {
		fmt.Fprintf(stderr, "scan failed: %v\n", err)
		return apperrors.ExitCode(err)
	}

	if err := writeScanResult(stdout, cfg.Format, result); err != nil {
		fmt.Fprintf(stderr, "report failed: %v\n", err)
		return apperrors.ExitCode(err)
	}

	return apperrors.ExitSuccess
}

func scanWithConfig(ctx context.Context, cfg config.ScanConfig) (report.ScanResult, error) {
	return newScanner(cfg).Scan(ctx, domain.ScanRequest{
		Inventory: domain.InventoryRequest{
			SteamID:   cfg.SteamID,
			Game:      cfg.Game,
			AppID:     cfg.AppID,
			ContextID: cfg.ContextID,
		},
		Prices: domain.PriceRequest{
			Path:     cfg.PriceCache,
			Currency: cfg.Currency,
		},
		Currency:       cfg.Currency,
		ThresholdCents: cfg.MinPriceCents,
		FeeBasisPoints: cfg.FeeBasisPoints,
	})
}

func writeScanResult(stdout io.Writer, format string, result report.ScanResult) error {
	switch format {
	case config.FormatJSON:
		return report.WriteJSON(stdout, result)
	case config.FormatTable:
		return report.WriteTable(stdout, result)
	default:
		return apperrors.New(apperrors.InvalidInput, "unsupported output format")
	}
}

func newScanner(cfg config.ScanConfig) app.Scanner {
	var inventoryProvider app.InventoryProvider
	if cfg.FixtureFile != "" {
		inventoryProvider = inventory.NewFixtureProvider(cfg.FixtureFile)
	} else {
		inventoryProvider = steam.NewInventoryProvider(steam.InventoryProviderOptions{
			Client: steam.NewClient(steam.Options{
				BaseURL: cfg.SteamBaseURL,
				Timeout: cfg.Timeout,
			}),
			CacheFile: cfg.CacheFile,
			CacheTTL:  cfg.CacheTTL,
			NoCache:   cfg.NoCache,
		})
	}
	return app.Scanner{
		InventoryProvider: inventoryProvider,
		PriceProvider:     pricing.FilePriceProvider{},
		Scorer:            pricing.Scorer{},
		Clock:             func() time.Time { return time.Now().UTC() },
	}
}

func runPrices(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: velesmist prices <template>")
		return apperrors.ExitInvalidInput
	}
	switch args[0] {
	case "template":
		cfg, err := config.ParsePriceTemplate(args[1:])
		if err != nil {
			fmt.Fprintf(stderr, "invalid input: %v\n", err)
			return apperrors.ExitCode(err)
		}
		if err := writePriceTemplate(cfg, stdout); err != nil {
			fmt.Fprintf(stderr, "prices template failed: %v\n", err)
			return apperrors.ExitCode(err)
		}
		return apperrors.ExitSuccess
	default:
		fmt.Fprintf(stderr, "unknown prices command: %s\n", args[0])
		return apperrors.ExitInvalidInput
	}
}

func writePriceTemplate(cfg config.PriceTemplateConfig, stdout io.Writer) error {
	if cfg.Output == "" {
		return pricing.WritePriceCacheTemplate(stdout)
	}
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if !cfg.Force {
		flags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	}
	file, err := os.OpenFile(cfg.Output, flags, 0o644)
	if err != nil {
		return apperrors.Wrap(apperrors.InvalidInput, "open price template output", err)
	}
	defer file.Close()
	return pricing.WritePriceCacheTemplate(file)
}
