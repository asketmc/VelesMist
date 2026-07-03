// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package config

import "testing"

func TestParseScanValidatesRequiredSteamID(t *testing.T) {
	_, err := ParseScan([]string{"--format", "json"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseScanAcceptsMinimumCommand(t *testing.T) {
	cfg, err := ParseScan([]string{"--steam-id", "76561198000000000", "--game", "dota2", "--format", "json", "--min-price", "5.00"})
	if err != nil {
		t.Fatalf("ParseScan returned error: %v", err)
	}
	if cfg.Game != GameDota2 || cfg.AppID != 570 || cfg.ContextID != "2" {
		t.Fatalf("unexpected Dota defaults: game=%s appid=%d contextid=%s", cfg.Game, cfg.AppID, cfg.ContextID)
	}
	if cfg.MinPriceCents != 500 {
		t.Fatalf("min price cents = %d, want 500", cfg.MinPriceCents)
	}
}

func TestParseScanAcceptsFixtureWithoutSteamID(t *testing.T) {
	cfg, err := ParseScan([]string{"--fixture", "testdata/dota_inventory.json", "--format", "json"})
	if err != nil {
		t.Fatalf("ParseScan returned error: %v", err)
	}
	if cfg.FixtureFile != "testdata/dota_inventory.json" {
		t.Fatalf("fixture = %s", cfg.FixtureFile)
	}
}

func TestValidateRejectsUnknownGame(t *testing.T) {
	cfg := defaultScanConfig()
	cfg.SteamID = "76561198000000000"
	cfg.Game = "tf2"
	if err := ValidateScan(cfg); err == nil {
		t.Fatal("expected invalid game error")
	}
}

func TestValidateRejectsInvalidFormat(t *testing.T) {
	cfg := defaultScanConfig()
	cfg.SteamID = "76561198000000000"
	cfg.Format = "xml"
	if err := ValidateScan(cfg); err == nil {
		t.Fatal("expected invalid format error")
	}
}

func TestParsePriceTemplateRejectsPositionalArgs(t *testing.T) {
	if _, err := ParsePriceTemplate([]string{"extra"}); err == nil {
		t.Fatal("expected positional argument error")
	}
}

func TestParsePriceTemplateAcceptsOutputAndForce(t *testing.T) {
	cfg, err := ParsePriceTemplate([]string{"--output", "price-cache.json", "--force"})
	if err != nil {
		t.Fatalf("ParsePriceTemplate error: %v", err)
	}
	if cfg.Output != "price-cache.json" || !cfg.Force {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
