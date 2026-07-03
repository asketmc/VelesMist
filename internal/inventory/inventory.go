// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package inventory

import (
	"encoding/json"
	"fmt"
	"sort"

	apperrors "github.com/asketmc/VelesMist/internal/errors"
)

type SteamInventory struct {
	Success      json.RawMessage    `json:"success"`
	Assets       []SteamAsset       `json:"assets"`
	Descriptions []SteamDescription `json:"descriptions"`
	TotalCount   int                `json:"total_inventory_count,omitempty"`
	More         bool               `json:"more_items,omitempty"`
	LastAssetID  string             `json:"last_assetid,omitempty"`
}

type SteamAsset struct {
	AppID      int    `json:"appid"`
	ContextID  string `json:"contextid"`
	AssetID    string `json:"assetid"`
	ClassID    string `json:"classid"`
	InstanceID string `json:"instanceid"`
	Amount     string `json:"amount"`
}

type SteamDescription struct {
	AppID          int    `json:"appid"`
	ClassID        string `json:"classid"`
	InstanceID     string `json:"instanceid"`
	Name           string `json:"name"`
	MarketHashName string `json:"market_hash_name"`
	Marketable     int    `json:"marketable"`
	Tradable       int    `json:"tradable"`
	Type           string `json:"type,omitempty"`
}

type Item struct {
	AssetID        string `json:"asset_id"`
	AppID          int    `json:"appid"`
	ClassID        string `json:"classid"`
	InstanceID     string `json:"instanceid"`
	Name           string `json:"name"`
	MarketHashName string `json:"market_hash_name"`
	Amount         int    `json:"amount"`
	Marketable     bool   `json:"marketable"`
	Tradable       bool   `json:"tradable"`
	Type           string `json:"type,omitempty"`
}

type AggregatedItem struct {
	AppID          int    `json:"appid"`
	Name           string `json:"name"`
	MarketHashName string `json:"market_hash_name"`
	Count          int    `json:"count"`
	Tradable       bool   `json:"tradable"`
}

func ParseSteamInventory(payload []byte) ([]Item, error) {
	var steamInv SteamInventory
	if err := json.Unmarshal(payload, &steamInv); err != nil {
		return nil, apperrors.Wrap(apperrors.InvalidInput, "decode Steam inventory JSON", err)
	}
	if !successOK(steamInv.Success) {
		return nil, apperrors.New(apperrors.Upstream, "Steam inventory response was not successful")
	}

	descriptions := make(map[string]SteamDescription, len(steamInv.Descriptions))
	for _, desc := range steamInv.Descriptions {
		descriptions[descriptionKey(desc.ClassID, desc.InstanceID)] = desc
	}

	items := make([]Item, 0, len(steamInv.Assets))
	for _, asset := range steamInv.Assets {
		desc, ok := descriptions[descriptionKey(asset.ClassID, asset.InstanceID)]
		if !ok {
			continue
		}
		name := desc.Name
		marketHashName := desc.MarketHashName
		if marketHashName == "" {
			marketHashName = name
		}
		if name == "" || marketHashName == "" {
			continue
		}
		amount := parseAmount(asset.Amount)
		items = append(items, Item{
			AssetID:        asset.AssetID,
			AppID:          firstPositive(asset.AppID, desc.AppID),
			ClassID:        asset.ClassID,
			InstanceID:     asset.InstanceID,
			Name:           name,
			MarketHashName: marketHashName,
			Amount:         amount,
			Marketable:     desc.Marketable == 1,
			Tradable:       desc.Tradable == 1,
			Type:           desc.Type,
		})
	}
	return items, nil
}

func AggregateMarketable(items []Item) []AggregatedItem {
	grouped := make(map[string]AggregatedItem)
	for _, item := range items {
		if !item.Marketable {
			continue
		}
		existing, ok := grouped[item.MarketHashName]
		if !ok {
			grouped[item.MarketHashName] = AggregatedItem{
				AppID:          item.AppID,
				Name:           item.Name,
				MarketHashName: item.MarketHashName,
				Count:          item.Amount,
				Tradable:       item.Tradable,
			}
			continue
		}
		existing.Count += item.Amount
		existing.Tradable = existing.Tradable && item.Tradable
		grouped[item.MarketHashName] = existing
	}

	rows := make([]AggregatedItem, 0, len(grouped))
	for _, item := range grouped {
		rows = append(rows, item)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].MarketHashName < rows[j].MarketHashName
	})
	return rows
}

func descriptionKey(classID string, instanceID string) string {
	if instanceID == "" {
		instanceID = "0"
	}
	return fmt.Sprintf("%s/%s", classID, instanceID)
}

func parseAmount(value string) int {
	if value == "" {
		return 1
	}
	var amount int
	if _, err := fmt.Sscanf(value, "%d", &amount); err != nil || amount <= 0 {
		return 1
	}
	return amount
}

func firstPositive(values ...int) int {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func successOK(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return true
	}
	var numeric int
	if err := json.Unmarshal(raw, &numeric); err == nil {
		return numeric != 0
	}
	var boolean bool
	if err := json.Unmarshal(raw, &boolean); err == nil {
		return boolean
	}
	return false
}
