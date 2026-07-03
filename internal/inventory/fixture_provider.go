// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package inventory

import (
	"context"
	"os"

	"github.com/asketmc/VelesMist/internal/domain"
	apperrors "github.com/asketmc/VelesMist/internal/errors"
)

type FixtureProvider struct {
	Path string
}

func NewFixtureProvider(path string) FixtureProvider {
	return FixtureProvider{Path: path}
}

func (p FixtureProvider) FetchInventory(_ context.Context, _ domain.InventoryRequest) (Inventory, error) {
	if p.Path == "" {
		return Inventory{}, apperrors.New(apperrors.InvalidInput, "fixture path is required")
	}
	body, err := os.ReadFile(p.Path)
	if err != nil {
		return Inventory{}, apperrors.Wrap(apperrors.InvalidInput, "read inventory fixture", err)
	}
	items, err := ParseSteamInventory(body)
	if err != nil {
		return Inventory{}, err
	}
	return Inventory{Items: items}, nil
}
