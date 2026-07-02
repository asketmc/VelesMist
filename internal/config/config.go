// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package config

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	apperrors "github.com/asketmc/VelesMist/internal/errors"
	"github.com/asketmc/VelesMist/internal/pricing"
)

const (
	FormatTable = "table"
	FormatJSON  = "json"
)

type ScanConfig struct {
	SteamID        string
	AppID          int
	ContextID      string
	Format         string
	Timeout        time.Duration
	CacheFile      string
	CacheTTL       time.Duration
	NoCache        bool
	PriceCache     string
	MinPriceCents  int64
	FeeBasisPoints int64
	Currency       string
	SteamBaseURL   string
}

func ParseScan(args []string) (ScanConfig, error) {
	cfg := defaultScanConfig()
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&cfg.SteamID, "steam-id", cfg.SteamID, "Steam64 ID for a public inventory")
	fs.IntVar(&cfg.AppID, "appid", cfg.AppID, "Steam appid")
	fs.StringVar(&cfg.ContextID, "contextid", cfg.ContextID, "Steam inventory contextid")
	fs.StringVar(&cfg.Format, "format", cfg.Format, "output format: table or json")
	fs.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "HTTP timeout")
	fs.StringVar(&cfg.CacheFile, "cache-file", cfg.CacheFile, "local JSON cache file")
	fs.DurationVar(&cfg.CacheTTL, "cache-ttl", cfg.CacheTTL, "inventory cache TTL")
	fs.BoolVar(&cfg.NoCache, "no-cache", cfg.NoCache, "disable local inventory cache")
	fs.StringVar(&cfg.PriceCache, "price-cache", cfg.PriceCache, "local JSON price cache")
	minPrice := fs.String("min-price", "0", "minimum estimated seller proceeds per item")
	fs.Int64Var(&cfg.FeeBasisPoints, "fee-bps", cfg.FeeBasisPoints, "estimated market fee in basis points")
	fs.StringVar(&cfg.Currency, "currency", cfg.Currency, "reporting currency label")
	fs.StringVar(&cfg.SteamBaseURL, "steam-base-url", cfg.SteamBaseURL, "Steam Community base URL")

	if err := fs.Parse(args); err != nil {
		return ScanConfig{}, apperrors.Wrap(apperrors.InvalidInput, "parse scan flags", err)
	}
	cents, err := pricing.ParseMoneyToCents(*minPrice)
	if err != nil {
		return ScanConfig{}, apperrors.Wrap(apperrors.InvalidInput, "parse min-price", err)
	}
	cfg.MinPriceCents = cents
	return cfg, ValidateScan(cfg)
}

func ValidateScan(cfg ScanConfig) error {
	if cfg.SteamID == "" {
		return apperrors.New(apperrors.InvalidInput, "steam-id is required")
	}
	if !isDigits(cfg.SteamID) || len(cfg.SteamID) < 16 || len(cfg.SteamID) > 20 {
		return apperrors.New(apperrors.InvalidInput, "steam-id must be a numeric Steam64 ID")
	}
	if cfg.AppID <= 0 {
		return apperrors.New(apperrors.InvalidInput, "appid must be positive")
	}
	if cfg.ContextID == "" || !isDigits(cfg.ContextID) {
		return apperrors.New(apperrors.InvalidInput, "contextid must be numeric")
	}
	if cfg.Format != FormatTable && cfg.Format != FormatJSON {
		return apperrors.New(apperrors.InvalidInput, "format must be table or json")
	}
	if cfg.Timeout <= 0 {
		return apperrors.New(apperrors.InvalidInput, "timeout must be positive")
	}
	if cfg.CacheTTL < 0 {
		return apperrors.New(apperrors.InvalidInput, "cache-ttl must be zero or positive")
	}
	if cfg.FeeBasisPoints < 0 {
		return apperrors.New(apperrors.InvalidInput, "fee-bps must be zero or positive")
	}
	if _, err := url.ParseRequestURI(cfg.SteamBaseURL); err != nil {
		return apperrors.Wrap(apperrors.InvalidInput, "steam-base-url is invalid", err)
	}
	return nil
}

func defaultScanConfig() ScanConfig {
	return ScanConfig{
		AppID:          570,
		ContextID:      "2",
		Format:         FormatTable,
		Timeout:        15 * time.Second,
		CacheFile:      defaultCacheFile(),
		CacheTTL:       10 * time.Minute,
		FeeBasisPoints: 1500,
		Currency:       "USD",
		SteamBaseURL:   "https://steamcommunity.com",
	}
}

func defaultCacheFile() string {
	if dir, err := os.UserCacheDir(); err == nil && dir != "" {
		return filepath.Join(dir, "velesmist", "cache.json")
	}
	return filepath.Join(".", ".velesmist-cache.json")
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	_, err := strconv.ParseUint(value, 10, 64)
	return err == nil
}

func Usage() string {
	return fmt.Sprintf("velesmist scan --steam-id <id> [--format %s|%s]", FormatTable, FormatJSON)
}
