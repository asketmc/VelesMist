// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package contracts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/asketmc/VelesMist/internal/pricing"
)

func TestSchemaFilesAreJSON(t *testing.T) {
	for _, path := range []string{
		"schemas/scan-report.v1.json",
		"schemas/price-cache.v1.json",
	} {
		var schema map[string]any
		readJSON(t, path, &schema)
		if stringField(schema, "$schema") == "" || stringField(schema, "$id") == "" {
			t.Fatalf("%s missing $schema or $id", path)
		}
	}
}

func TestScanReportGoldenMatchesContract(t *testing.T) {
	var report map[string]any
	readJSON(t, "internal/report/testdata/scan.json.golden", &report)
	if err := validateScanReport(report); err != nil {
		t.Fatalf("scan report contract violation: %v", err)
	}
}

func TestScanReportInvalidFixtureFailsContract(t *testing.T) {
	var report map[string]any
	readJSON(t, "schemas/testdata/scan-report.invalid-recommendation.json", &report)
	if err := validateScanReport(report); err == nil {
		t.Fatal("expected invalid scan report fixture to fail")
	}
}

func TestPriceCacheFixtureMatchesContractAndLoader(t *testing.T) {
	body := readFile(t, "schemas/testdata/price-cache.valid.json")
	var cache map[string]any
	decodeJSON(t, "schemas/testdata/price-cache.valid.json", body, &cache)
	if err := validatePriceCache(cache); err != nil {
		t.Fatalf("price cache contract violation: %v", err)
	}
	prices, meta, err := pricing.LoadPriceCache(bytes.NewReader(body))
	if err != nil {
		t.Fatalf("valid price cache fixture did not load: %v", err)
	}
	if meta.SchemaVersion != pricing.PriceCacheSchemaVersion || meta.Count != 2 || len(prices) != 2 {
		t.Fatalf("unexpected loader result: meta=%+v prices=%+v", meta, prices)
	}
}

func TestPriceCacheInvalidFixturesFailContract(t *testing.T) {
	for _, path := range []string{
		"schemas/testdata/price-cache.invalid-schema.json",
		"schemas/testdata/price-cache.invalid-missing-fields.json",
		"schemas/testdata/price-cache.invalid-missing-price.json",
	} {
		var cache map[string]any
		readJSON(t, path, &cache)
		if err := validatePriceCache(cache); err == nil {
			t.Fatalf("%s should fail price cache contract", path)
		}
	}
}

func validateScanReport(report map[string]any) error {
	if err := rejectUnknownKeys(report,
		"schema_version",
		"generated_at",
		"steam_id",
		"appid",
		"contextid",
		"currency",
		"threshold_cents",
		"items",
		"candidates",
		"summary",
	); err != nil {
		return err
	}
	required := []string{
		"schema_version",
		"generated_at",
		"steam_id",
		"appid",
		"contextid",
		"currency",
		"threshold_cents",
		"items",
		"candidates",
		"summary",
	}
	if err := requireKeys(report, required...); err != nil {
		return err
	}
	if got := stringField(report, "schema_version"); got != "velesmist.scan.v1" {
		return fmt.Errorf("schema_version = %q", got)
	}
	if _, err := time.Parse(time.RFC3339, stringField(report, "generated_at")); err != nil {
		return fmt.Errorf("generated_at is not RFC3339: %w", err)
	}
	if stringField(report, "currency") == "" || numberField(report, "appid") <= 0 {
		return fmt.Errorf("invalid report target fields")
	}
	items, ok := report["items"].([]any)
	if !ok {
		return fmt.Errorf("items must be an array")
	}
	for i, raw := range items {
		item, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("items[%d] must be an object", i)
		}
		if err := validateScanItem(item); err != nil {
			return fmt.Errorf("items[%d]: %w", i, err)
		}
	}
	candidates, ok := report["candidates"].([]any)
	if !ok {
		return fmt.Errorf("candidates must be an array")
	}
	for i, raw := range candidates {
		item, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("candidates[%d] must be an object", i)
		}
		if err := validateScanItem(item); err != nil {
			return fmt.Errorf("candidates[%d]: %w", i, err)
		}
		if stringField(item, "recommendation") != "sell" || boolField(item, "candidate") != true {
			return fmt.Errorf("candidates[%d] must be sell/candidate=true", i)
		}
	}
	summary, ok := report["summary"].(map[string]any)
	if !ok {
		return fmt.Errorf("summary must be an object")
	}
	if err := rejectUnknownKeys(summary,
		"marketable_items",
		"priced_items",
		"missing_price_items",
		"skipped_items",
		"candidate_items",
		"estimated_total_gross_cents",
		"estimated_total_fee_cents",
		"estimated_total_receive_cents",
	); err != nil {
		return err
	}
	return requireNumberKeys(summary,
		"marketable_items",
		"priced_items",
		"missing_price_items",
		"skipped_items",
		"candidate_items",
		"estimated_total_gross_cents",
		"estimated_total_fee_cents",
		"estimated_total_receive_cents",
	)
}

func validateScanItem(item map[string]any) error {
	if err := rejectUnknownKeys(item,
		"appid",
		"name",
		"market_hash_name",
		"count",
		"tradable",
		"market_url",
		"price_status",
		"buyer_price_cents",
		"estimated_fee_cents",
		"seller_receive_cents",
		"total_buyer_price_cents",
		"total_estimated_fee_cents",
		"total_receive_cents",
		"price_source",
		"liquidity_score",
		"confidence",
		"recommendation",
		"reason_codes",
		"candidate",
	); err != nil {
		return err
	}
	if err := requireKeys(item,
		"appid",
		"name",
		"market_hash_name",
		"count",
		"tradable",
		"market_url",
		"price_status",
		"liquidity_score",
		"confidence",
		"recommendation",
		"reason_codes",
		"candidate",
	); err != nil {
		return err
	}
	if numberField(item, "appid") <= 0 || numberField(item, "count") <= 0 {
		return fmt.Errorf("appid/count must be positive")
	}
	if stringField(item, "name") == "" || stringField(item, "market_hash_name") == "" {
		return fmt.Errorf("name and market_hash_name are required")
	}
	if !strings.HasPrefix(stringField(item, "market_url"), "https://steamcommunity.com/market/listings/") {
		return fmt.Errorf("market_url must point to Steam Community market")
	}
	switch stringField(item, "price_status") {
	case "priced":
		if err := requireNumberKeys(item,
			"buyer_price_cents",
			"estimated_fee_cents",
			"seller_receive_cents",
			"total_buyer_price_cents",
			"total_estimated_fee_cents",
			"total_receive_cents",
		); err != nil {
			return err
		}
		if stringField(item, "price_source") == "" {
			return fmt.Errorf("priced items must include price_source")
		}
	case "missing":
		if err := rejectKeys(item,
			"buyer_price_cents",
			"estimated_fee_cents",
			"seller_receive_cents",
			"total_buyer_price_cents",
			"total_estimated_fee_cents",
			"total_receive_cents",
			"price_source",
		); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid price_status %q", stringField(item, "price_status"))
	}
	switch stringField(item, "recommendation") {
	case "sell", "skip", "missing_price":
	default:
		return fmt.Errorf("invalid recommendation %q", stringField(item, "recommendation"))
	}
	switch stringField(item, "confidence") {
	case "medium", "none":
	default:
		return fmt.Errorf("invalid confidence %q", stringField(item, "confidence"))
	}
	reasons, ok := item["reason_codes"].([]any)
	if !ok || len(reasons) == 0 {
		return fmt.Errorf("reason_codes must be a non-empty array")
	}
	for _, reason := range reasons {
		if text, ok := reason.(string); !ok || text == "" {
			return fmt.Errorf("reason_codes entries must be non-empty strings")
		}
	}
	if _, ok := item["tradable"].(bool); !ok {
		return fmt.Errorf("tradable must be boolean")
	}
	if _, ok := item["candidate"].(bool); !ok {
		return fmt.Errorf("candidate must be boolean")
	}
	switch stringField(item, "recommendation") {
	case "sell":
		if stringField(item, "price_status") != "priced" || !boolField(item, "candidate") {
			return fmt.Errorf("sell recommendation must be priced and candidate=true")
		}
	case "skip":
		if stringField(item, "price_status") != "priced" || boolField(item, "candidate") {
			return fmt.Errorf("skip recommendation must be priced and candidate=false")
		}
	case "missing_price":
		if stringField(item, "price_status") != "missing" || stringField(item, "confidence") != "none" || boolField(item, "candidate") {
			return fmt.Errorf("missing_price recommendation must be missing/confidence=none/candidate=false")
		}
	}
	return nil
}

func validatePriceCache(cache map[string]any) error {
	if err := rejectUnknownKeys(cache, "schema_version", "currency", "prices"); err != nil {
		return err
	}
	if err := requireKeys(cache, "schema_version", "currency", "prices"); err != nil {
		return err
	}
	if got := stringField(cache, "schema_version"); got != "velesmist.price-cache.v1" {
		return fmt.Errorf("schema_version = %q", got)
	}
	if stringField(cache, "currency") == "" {
		return fmt.Errorf("currency is required")
	}
	prices, ok := cache["prices"].(map[string]any)
	if !ok {
		return fmt.Errorf("prices must be an object")
	}
	for marketHashName, raw := range prices {
		if marketHashName == "" {
			return fmt.Errorf("prices key must be market_hash_name")
		}
		price, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("price for %s must be an object", marketHashName)
		}
		if err := validatePriceEntry(price); err != nil {
			return fmt.Errorf("price for %s: %w", marketHashName, err)
		}
	}
	return nil
}

func validatePriceEntry(price map[string]any) error {
	allowed := map[string]bool{
		"buyer_price_cents": true,
		"lowest_price":      true,
		"median_price":      true,
		"source":            true,
		"confidence":        true,
		"liquidity_score":   true,
	}
	for key := range price {
		if !allowed[key] {
			return fmt.Errorf("unknown field %q", key)
		}
	}
	hasPrice := false
	if raw, ok := price["buyer_price_cents"]; ok {
		value, ok := raw.(float64)
		if !ok || value <= 0 || value != float64(int64(value)) {
			return fmt.Errorf("buyer_price_cents must be a positive integer")
		}
		hasPrice = true
	}
	for _, key := range []string{"lowest_price", "median_price"} {
		if value := stringField(price, key); value != "" {
			hasPrice = true
		}
	}
	if confidence := stringField(price, "confidence"); confidence != "" {
		switch confidence {
		case "medium", "none":
		default:
			return fmt.Errorf("invalid confidence %q", confidence)
		}
	}
	if raw, ok := price["liquidity_score"]; ok {
		value, ok := raw.(float64)
		if !ok || value < 0 || value != float64(int64(value)) {
			return fmt.Errorf("liquidity_score must be a non-negative integer")
		}
	}
	if !hasPrice {
		return fmt.Errorf("one of buyer_price_cents, lowest_price, or median_price is required")
	}
	return nil
}

func requireKeys(object map[string]any, keys ...string) error {
	for _, key := range keys {
		if _, ok := object[key]; !ok {
			return fmt.Errorf("missing required key %q", key)
		}
	}
	return nil
}

func rejectUnknownKeys(object map[string]any, keys ...string) error {
	allowed := make(map[string]bool, len(keys))
	for _, key := range keys {
		allowed[key] = true
	}
	for key := range object {
		if !allowed[key] {
			return fmt.Errorf("unknown field %q", key)
		}
	}
	return nil
}

func rejectKeys(object map[string]any, keys ...string) error {
	for _, key := range keys {
		if _, ok := object[key]; ok {
			return fmt.Errorf("field %q is not allowed", key)
		}
	}
	return nil
}

func requireNumberKeys(object map[string]any, keys ...string) error {
	for _, key := range keys {
		raw, ok := object[key]
		if !ok {
			return fmt.Errorf("missing required number key %q", key)
		}
		value, ok := raw.(float64)
		if !ok || value < 0 || value != float64(int64(value)) {
			return fmt.Errorf("%s must be a non-negative integer", key)
		}
	}
	return nil
}

func stringField(object map[string]any, key string) string {
	value, _ := object[key].(string)
	return value
}

func numberField(object map[string]any, key string) float64 {
	value, _ := object[key].(float64)
	return value
}

func boolField(object map[string]any, key string) bool {
	value, _ := object[key].(bool)
	return value
}

func readJSON(t *testing.T, rel string, out any) {
	t.Helper()
	body := readFile(t, rel)
	decodeJSON(t, rel, body, out)
}

func readFile(t *testing.T, rel string) []byte {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(repoRoot(t), filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return body
}

func decodeJSON(t *testing.T, rel string, body []byte, out any) {
	t.Helper()
	if err := json.Unmarshal(body, out); err != nil {
		t.Fatalf("decode %s: %v", rel, err)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
