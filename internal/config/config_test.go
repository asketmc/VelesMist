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
	cfg, err := ParseScan([]string{"--steam-id", "76561198000000000", "--format", "json", "--min-price", "5.00"})
	if err != nil {
		t.Fatalf("ParseScan returned error: %v", err)
	}
	if cfg.AppID != 570 || cfg.ContextID != "2" {
		t.Fatalf("unexpected Dota defaults: appid=%d contextid=%s", cfg.AppID, cfg.ContextID)
	}
	if cfg.MinPriceCents != 500 {
		t.Fatalf("min price cents = %d, want 500", cfg.MinPriceCents)
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
