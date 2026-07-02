// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	apperrors "github.com/asketmc/VelesMist/internal/errors"
)

const SchemaVersion = 1

type Store struct {
	path string
}

type File struct {
	SchemaVersion int               `json:"schema_version"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Entries       map[string]Record `json:"entries"`
}

type Record struct {
	FetchedAt time.Time       `json:"fetched_at"`
	ExpiresAt time.Time       `json:"expires_at"`
	Body      json.RawMessage `json:"body"`
}

func NewStore(path string) Store {
	return Store{path: path}
}

func InventoryKey(steamID string, appID int, contextID string) string {
	return fmt.Sprintf("inventory:%s:%d:%s", steamID, appID, contextID)
}

func (s Store) GetValid(key string, now time.Time) ([]byte, bool, error) {
	file, err := s.read()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	record, ok := file.Entries[key]
	if !ok || now.After(record.ExpiresAt) {
		return nil, false, nil
	}
	return append([]byte(nil), record.Body...), true, nil
}

func (s Store) Put(key string, body []byte, now time.Time, ttl time.Duration) error {
	file, err := s.read()
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		file = File{SchemaVersion: SchemaVersion, Entries: map[string]Record{}}
	}
	if file.Entries == nil {
		file.Entries = map[string]Record{}
	}
	file.SchemaVersion = SchemaVersion
	file.UpdatedAt = now
	file.Entries[key] = Record{
		FetchedAt: now,
		ExpiresAt: now.Add(ttl),
		Body:      append([]byte(nil), body...),
	}
	return s.write(file)
}

func (s Store) read() (File, error) {
	body, err := os.ReadFile(s.path)
	if err != nil {
		return File{}, err
	}
	var file File
	if err := json.Unmarshal(body, &file); err != nil {
		return File{}, apperrors.Wrap(apperrors.InvalidInput, "decode cache file", err)
	}
	if file.SchemaVersion != SchemaVersion {
		return File{}, apperrors.New(apperrors.InvalidInput, "unsupported cache schema version")
	}
	if file.Entries == nil {
		file.Entries = map[string]Record{}
	}
	return file, nil
}

func (s Store) write(file File) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return apperrors.Wrap(apperrors.Internal, "create cache directory", err)
	}
	body, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return apperrors.Wrap(apperrors.Internal, "encode cache file", err)
	}
	if err := os.WriteFile(s.path, append(body, '\n'), 0o600); err != nil {
		return apperrors.Wrap(apperrors.Internal, "write cache file", err)
	}
	return nil
}
