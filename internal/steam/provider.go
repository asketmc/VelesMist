// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package steam

import (
	"context"
	"time"

	"github.com/asketmc/VelesMist/internal/cache"
	"github.com/asketmc/VelesMist/internal/domain"
	"github.com/asketmc/VelesMist/internal/inventory"
)

type InventoryProvider struct {
	client    Client
	cacheFile string
	cacheTTL  time.Duration
	noCache   bool
	clock     func() time.Time
}

type InventoryProviderOptions struct {
	Client    Client
	CacheFile string
	CacheTTL  time.Duration
	NoCache   bool
	Clock     func() time.Time
}

func NewInventoryProvider(opts InventoryProviderOptions) InventoryProvider {
	if opts.Clock == nil {
		opts.Clock = func() time.Time { return time.Now().UTC() }
	}
	return InventoryProvider{
		client:    opts.Client,
		cacheFile: opts.CacheFile,
		cacheTTL:  opts.CacheTTL,
		noCache:   opts.NoCache,
		clock:     opts.Clock,
	}
}

func (p InventoryProvider) FetchInventory(ctx context.Context, req domain.InventoryRequest) (inventory.Inventory, error) {
	cacheKey := cache.InventoryKey(req.SteamID, req.AppID, req.ContextID)
	store := cache.NewStore(p.cacheFile)
	now := p.clock().UTC()
	if !p.noCache {
		if body, ok, err := store.GetValid(cacheKey, now); err == nil && ok {
			return parseInventory(body)
		}
	}

	payload, err := p.client.FetchInventory(ctx, req.SteamID, req.AppID, req.ContextID)
	if err != nil {
		return inventory.Inventory{}, err
	}
	if !p.noCache {
		_ = store.Put(cacheKey, payload, now, p.cacheTTL)
	}
	return parseInventory(payload)
}

func parseInventory(payload []byte) (inventory.Inventory, error) {
	items, err := inventory.ParseSteamInventory(payload)
	if err != nil {
		return inventory.Inventory{}, err
	}
	return inventory.Inventory{Items: items}, nil
}
