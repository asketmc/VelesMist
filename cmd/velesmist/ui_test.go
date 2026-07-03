// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/asketmc/VelesMist/internal/config"
	apperrors "github.com/asketmc/VelesMist/internal/errors"
	"github.com/asketmc/VelesMist/internal/pricing"
	"github.com/asketmc/VelesMist/internal/report"
)

func TestRunDispatchesUICommandValidation(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"ui", "--addr", "0.0.0.0:8765"}, &stdout, &stderr)

	if code != apperrors.ExitInvalidInput {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "localhost") {
		t.Fatalf("stderr = %q, want localhost validation", stderr.String())
	}
}

func TestParseUIRejectsNonLocalhostAddress(t *testing.T) {
	_, err := parseUI([]string{"--addr", "0.0.0.0:8765"})
	if err == nil {
		t.Fatal("expected non-localhost bind to fail")
	}
}

func TestParseUIContract(t *testing.T) {
	cfg, err := parseUI([]string{"--addr", "localhost:0", "--open=false"})
	if err != nil {
		t.Fatalf("parseUI returned error: %v", err)
	}
	if cfg.Addr != "localhost:0" || cfg.Open {
		t.Fatalf("unexpected config: %+v", cfg)
	}

	if _, err := parseUI([]string{"extra"}); err == nil {
		t.Fatal("expected positional argument error")
	}
	if _, err := parseUI([]string{"--addr", "not-a-host-port"}); err == nil {
		t.Fatal("expected malformed address error")
	}
}

func TestServeUIContract(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	openCh := make(chan string, 1)
	codeCh := make(chan int, 1)
	go func() {
		codeCh <- serveUI(ctx, listener, newUIHandler(func(context.Context, config.ScanConfig) (report.ScanResult, error) {
			return report.ScanResult{}, nil
		}), true, func(url string) error {
			openCh <- url
			return nil
		}, &stdout, &stderr)
	}()

	url := waitForOpenURL(t, openCh)
	if !strings.HasPrefix(url, "http://127.0.0.1:") {
		t.Fatalf("open URL = %s", url)
	}
	waitForHealthz(t, url+"healthz")
	cancel()
	select {
	case code := <-codeCh:
		if code != apperrors.ExitSuccess {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serveUI did not stop after context cancellation")
	}
	if !strings.Contains(stdout.String(), "VelesMist UI listening on "+url) {
		t.Fatalf("stdout = %q, want listening URL %s", stdout.String(), url)
	}
}

func TestServeUIReportsBrowserOpenFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	openCh := make(chan string, 1)
	codeCh := make(chan int, 1)
	go func() {
		codeCh <- serveUI(ctx, listener, newUIHandler(func(context.Context, config.ScanConfig) (report.ScanResult, error) {
			return report.ScanResult{}, nil
		}), true, func(url string) error {
			openCh <- url
			return errors.New("browser unavailable")
		}, &stdout, &stderr)
	}()

	url := waitForOpenURL(t, openCh)
	waitForHealthz(t, url+"healthz")
	cancel()
	select {
	case code := <-codeCh:
		if code != apperrors.ExitSuccess {
			t.Fatalf("exit=%d", code)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serveUI did not stop")
	}
	if !strings.Contains(stderr.String(), "browser open failed") || !strings.Contains(stdout.String(), "VelesMist UI listening on "+url) {
		t.Fatalf("stdout=%q stderr=%q, want listening URL and browser open failure", stdout.String(), stderr.String())
	}
}

func TestServeUIReportsServeFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := serveUI(context.Background(), listener, newUIHandler(func(context.Context, config.ScanConfig) (report.ScanResult, error) {
		return report.ScanResult{}, nil
	}), false, func(string) error {
		t.Fatal("open should not be called")
		return nil
	}, &stdout, &stderr)

	if code != apperrors.ExitInternal {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "ui failed") {
		t.Fatalf("stderr = %q, want serve failure", stderr.String())
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

func TestUIHTMLContract(t *testing.T) {
	required := []string{
		`<form id="scan-form">`,
		`name="steam_id"`,
		`name="game"`,
		`name="price_cache"`,
		`name="fixture"`,
		`name="no_cache"`,
		`fetch('/api/scan'`,
		`method: 'POST'`,
		`Table`,
		`JSON`,
		`No login, cookies, listings, or selling.`,
	}
	for _, text := range required {
		if !strings.Contains(uiHTML, text) {
			t.Fatalf("uiHTML missing contract text %q", text)
		}
	}
	forbidden := []string{
		`<script src=`,
		`<link`,
		`https://`,
		`http://`,
	}
	for _, text := range forbidden {
		if strings.Contains(uiHTML, text) {
			t.Fatalf("uiHTML contains forbidden external surface %q", text)
		}
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

func TestUIScanSuccessContractWithMockScanner(t *testing.T) {
	wantReport := report.ScanResult{
		SchemaVersion:  report.SchemaVersion,
		SteamID:        "76561198000000000",
		AppID:          570,
		ContextID:      "2",
		Currency:       "EUR",
		ThresholdCents: 500,
		Items: []pricing.PricedItem{{
			AppID:          570,
			Name:           "Golden Moonfall",
			MarketHashName: "Golden Moonfall",
			Count:          1,
			MarketURL:      "https://steamcommunity.com/market/listings/570/Golden%20Moonfall",
			PriceStatus:    pricing.PriceStatusPriced,
			Recommendation: pricing.RecommendationSell,
			ReasonCodes:    []string{pricing.ReasonMarketable, pricing.ReasonPriceFound},
			Candidate:      true,
		}},
		Summary: report.Summary{CandidateItems: 1, PricedItems: 1},
	}
	handler := newUIHandler(func(ctx context.Context, cfg config.ScanConfig) (report.ScanResult, error) {
		if cfg.SteamID != "76561198000000000" || cfg.Game != config.GameDota2 || cfg.Format != config.FormatJSON {
			t.Fatalf("unexpected scan config: %+v", cfg)
		}
		if cfg.PriceCache != "prices.json" || cfg.MinPriceCents != 500 || cfg.Currency != "EUR" || !cfg.NoCache {
			t.Fatalf("unexpected optional scan config: %+v", cfg)
		}
		return wantReport, nil
	})
	req := httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(`{
		"steam_id":" 76561198000000000 ",
		"game":"dota2",
		"price_cache":" prices.json ",
		"min_price":"5.00",
		"currency":"EUR",
		"no_cache":true
	}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("content-type = %q", got)
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("cache-control = %q", got)
	}
	var body uiScanResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.OK || body.Report.SchemaVersion != report.SchemaVersion || body.Report.Items[0].Recommendation != pricing.RecommendationSell {
		t.Fatalf("unexpected response: %+v", body)
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

func TestUIScanRejectsUnknownFieldAndUnsupportedGame(t *testing.T) {
	handler := newUIHandler(scanWithConfig)
	tests := []string{
		`{"steam_id":"76561198000000000","unknown":true}`,
		`{"steam_id":"76561198000000000","game":"tf2"}`,
	}
	for _, payload := range tests {
		req := httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(payload))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("payload=%s status=%d body=%s", payload, rec.Code, rec.Body.String())
		}
		var got uiErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if got.OK || got.Kind != string(apperrors.InvalidInput) || got.ExitCode != apperrors.ExitInvalidInput {
			t.Fatalf("payload=%s unexpected response: %+v", payload, got)
		}
	}
}

func TestUIScanMethodContract(t *testing.T) {
	handler := newUIHandler(scanWithConfig)
	req := httptest.NewRequest(http.MethodGet, "/api/scan", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Allow"); got != http.MethodPost {
		t.Fatalf("Allow = %q", got)
	}
	var body uiErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.OK || body.Kind != string(apperrors.InvalidInput) {
		t.Fatalf("unexpected response: %+v", body)
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

func TestUIErrorContractByKind(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		kind   apperrors.Kind
		text   string
	}{
		{
			name:   "rate limit",
			err:    apperrors.New(apperrors.RateLimited, "raw rate details"),
			status: http.StatusTooManyRequests,
			kind:   apperrors.RateLimited,
			text:   "rate limited",
		},
		{
			name:   "timeout",
			err:    apperrors.New(apperrors.NetworkTimeout, "raw timeout details"),
			status: http.StatusGatewayTimeout,
			kind:   apperrors.NetworkTimeout,
			text:   "timed out",
		},
		{
			name:   "internal",
			err:    errors.New("database exploded with secret"),
			status: http.StatusInternalServerError,
			kind:   apperrors.Internal,
			text:   "Internal error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeUIError(rec, tt.err)
			if rec.Code != tt.status {
				t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
			}
			var got uiErrorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if got.OK || got.Kind != string(tt.kind) || got.ExitCode != apperrors.ExitCode(tt.err) {
				t.Fatalf("unexpected response: %+v", got)
			}
			if !strings.Contains(got.Message, tt.text) {
				t.Fatalf("message = %q, want %q", got.Message, tt.text)
			}
			if strings.Contains(got.Message, "secret") || strings.Contains(got.Message, "raw") {
				t.Fatalf("message leaked raw error: %q", got.Message)
			}
		})
	}
}

func TestScanConfigFromUIContract(t *testing.T) {
	cfg, err := scanConfigFromUI(uiScanRequest{
		SteamID:    " 76561198000000000 ",
		Game:       "",
		PriceCache: " prices.json ",
		MinPrice:   "7.25",
		Currency:   "EUR",
		NoCache:    true,
	})
	if err != nil {
		t.Fatalf("scanConfigFromUI error: %v", err)
	}
	if cfg.SteamID != "76561198000000000" || cfg.Game != config.GameDota2 || cfg.Format != config.FormatJSON {
		t.Fatalf("unexpected required config: %+v", cfg)
	}
	if cfg.PriceCache != "prices.json" || cfg.MinPriceCents != 725 || cfg.Currency != "EUR" || !cfg.NoCache {
		t.Fatalf("unexpected optional config: %+v", cfg)
	}
}

func TestScanConfigFromUIFixtureContract(t *testing.T) {
	cfg, err := scanConfigFromUI(uiScanRequest{
		SteamID: "not-used-when-fixture-is-set",
		Fixture: " inventory.json ",
	})
	if err != nil {
		t.Fatalf("scanConfigFromUI error: %v", err)
	}
	if cfg.FixtureFile != "inventory.json" || cfg.SteamID != "" || cfg.Format != config.FormatJSON {
		t.Fatalf("unexpected fixture config: %+v", cfg)
	}
}

func waitForOpenURL(t *testing.T, ch <-chan string) string {
	t.Helper()
	select {
	case url := <-ch:
		return url
	case <-time.After(5 * time.Second):
		t.Fatal("browser open callback was not called")
		return ""
	}
}

func waitForHealthz(t *testing.T, url string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			body := new(bytes.Buffer)
			_, _ = body.ReadFrom(resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK && body.String() == "ok\n" {
				return
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("health endpoint did not become ready: %s", url)
}
