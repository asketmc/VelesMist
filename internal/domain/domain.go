// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package domain

type InventoryRequest struct {
	SteamID   string
	Game      string
	AppID     int
	ContextID string
}

type PriceRequest struct {
	Path     string
	Currency string
}

type ScanRequest struct {
	Inventory      InventoryRequest
	Prices         PriceRequest
	Currency       string
	ThresholdCents int64
	FeeBasisPoints int64
}
