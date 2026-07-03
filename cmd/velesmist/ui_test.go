// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asketmc/VelesMist/internal/config"
	apperrors "github.com/asketmc/VelesMist/internal/errors"
	"github.com/asketmc/VelesMist/internal/report"
)

func TestParseUIRejectsNonLocalhostAddress(t *testing.T) {
	_, err := parseUI([]string{"--addr", "0.0.0.0:8765"})
	if err == nil {
		t.Fatal("expected non-localhost bind to fail")
	}
}

func TestUIIndexServesMinimalPage(t *testing.T) {
	handler := newUIHandler(func(context.Context, config.ScanConfig) (report.ScanResult, error) {
		t.Fatal("scan should not run for index")
		return report.ScanResult{}, nil
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "VelesMist") || !strings.Contains(body, "No login") {
		t.Fatalf("unexpected UI body: %s", body)
	}
	if got := rec.Header().Get("Content-Security-Policy"); !strings.Contains(got, "default-src 'none'") {
		t.Fatalf("missing strict CSP: %s", got)
	}
}

func TestUIScanUsesFixtureWithoutNetwork(t *testing.T) {
	handler := newUIHandler(scanWithConfig)
	payload := map[string]any{
		"fixture": filepath.Join("..", "..", "internal", "inventory", "testdata", "dota_inventory.json"),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/scan", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var got uiScanResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v\n%s", err, rec.Body.String())
	}
	if !got.OK || got.Report.SchemaVersion != report.SchemaVersion {
		t.Fatalf("unexpected response: %+v", got)
	}
	if got.Report.Summary.MarketableItems != 2 || got.Report.Summary.MissingPriceItems != 2 {
		t.Fatalf("unexpected summary: %+v", got.Report.Summary)
	}
}

func TestUIScanRejectsMissingSteamID(t *testing.T) {
	handler := newUIHandler(scanWithConfig)
	req := httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(`{"game":"dota2"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var got uiErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.OK || got.Kind != string(apperrors.InvalidInput) || got.ExitCode != apperrors.ExitInvalidInput {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestUIScanSanitizesUpstreamError(t *testing.T) {
	handler := newUIHandler(func(context.Context, config.ScanConfig) (report.ScanResult, error) {
		return report.ScanResult{}, apperrors.Wrap(apperrors.Upstream, "Steam request failed", context.DeadlineExceeded)
	})
	req := httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(`{"steam_id":"76561197987179126","game":"dota2"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "76561197987179126") || strings.Contains(rec.Body.String(), "steamcommunity.com") {
		t.Fatalf("response leaked request details: %s", rec.Body.String())
	}
	var got uiErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Kind != string(apperrors.Upstream) || got.ExitCode != apperrors.ExitUpstream {
		t.Fatalf("unexpected error response: %+v", got)
	}
}
