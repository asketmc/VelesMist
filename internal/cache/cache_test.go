// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package cache

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreReadWriteValidRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	store := NewStore(path)
	now := time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC)
	key := InventoryKey("76561198000000000", 570, "2")
	body := []byte(`{"success":1}`)
	if err := store.Put(key, body, now, time.Minute); err != nil {
		t.Fatalf("Put error: %v", err)
	}
	got, ok, err := store.GetValid(key, now.Add(30*time.Second))
	if err != nil {
		t.Fatalf("GetValid error: %v", err)
	}
	if !ok {
		t.Fatal("expected cache hit")
	}
	if !json.Valid(got) {
		t.Fatalf("cached body is not valid JSON: %s", got)
	}
	if !json.Valid(body) {
		t.Fatalf("test body is not valid JSON: %s", body)
	}
	if !jsonEqual(got, body) {
		t.Fatalf("body = %s, want %s", got, body)
	}
}

func TestStoreExpiresRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	store := NewStore(path)
	now := time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC)
	key := InventoryKey("76561198000000000", 570, "2")
	if err := store.Put(key, []byte(`{}`), now, time.Second); err != nil {
		t.Fatalf("Put error: %v", err)
	}
	_, ok, err := store.GetValid(key, now.Add(2*time.Second))
	if err != nil {
		t.Fatalf("GetValid error: %v", err)
	}
	if ok {
		t.Fatal("expected expired cache miss")
	}
}

func jsonEqual(a, b []byte) bool {
	var left any
	var right any
	if err := json.Unmarshal(a, &left); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &right); err != nil {
		return false
	}
	leftBytes, err := json.Marshal(left)
	if err != nil {
		return false
	}
	rightBytes, err := json.Marshal(right)
	if err != nil {
		return false
	}
	return string(leftBytes) == string(rightBytes)
}
