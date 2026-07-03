# Scan Report Contract v1

Machine-readable schema: [`schemas/scan-report.v1.json`](../../schemas/scan-report.v1.json).

The JSON report contract is stable for automation within schema version `velesmist.scan.v1`.

## Required Top-Level Fields

- `schema_version`: must be `velesmist.scan.v1`.
- `generated_at`: RFC3339 UTC timestamp.
- `steam_id`: Steam64 ID when scanning a real inventory; may be empty for fixture-only scans.
- `appid`: Steam app ID, currently `570` for Dota 2.
- `contextid`: Steam inventory context ID, currently `2` for Dota 2.
- `currency`: report currency label, default `USD`.
- `threshold_cents`: minimum seller receive threshold in cents.
- `items`: all marketable items that the current scorer evaluates.
- `candidates`: subset of `items` with `recommendation: "sell"`.
- `summary`: aggregate counts and estimated totals.

## Item Fields

- Item identity fields: `appid`, `name`, `market_hash_name`, `count`, `tradable`, and `market_url`.
- `market_hash_name`: Steam market hash name used as the price-cache lookup key and market URL identity.
- `recommendation`: one of `sell`, `skip`, `missing_price`.
- `not_marketable`: planned recommendation value for a future product PR; it is not emitted by the current v1 implementation.
- `price_status`: one of `priced`, `missing`.
- `buyer_price_cents`: gross buyer-facing price when priced.
- `estimated_fee_cents`: estimated Steam/platform fee for one item when priced.
- `seller_receive_cents`: estimated seller receive for one item when priced.
- `total_buyer_price_cents`: gross buyer-facing price multiplied by count when priced.
- `total_estimated_fee_cents`: estimated fee multiplied by count when priced.
- `total_receive_cents`: estimated seller receive multiplied by count when priced.
- `confidence`: `medium` for local cache prices, `none` when price is missing.
- `liquidity_score`: integer score; v1 uses `0` because live liquidity is not implemented.
- `reason_codes`: non-empty explanation list.
- `market_url`: Steam Community Market URL derived from `market_hash_name`.
- `candidate`: true only for `sell` recommendations.

Priced items must include gross buyer price, estimated fee, seller receive, total gross buyer price, total estimated fee, total receive, and `price_source`.

Missing-price items use `price_status: "missing"`, `recommendation: "missing_price"`, `confidence: "none"`, `candidate: false`, and omit gross/fee/receive fields.

## Currency Behavior

The report `currency` field is a label inherited from config or the local price cache. VelesMist does not convert currencies in `velesmist.scan.v1`; all `*_cents` values are interpreted in that selected currency.

## Sorting Rules

Report rows are sorted by recommendation priority, then total seller receive, buyer price, count, and `market_hash_name`.

Recommendation priority:

1. `sell`
2. `skip`
3. `missing_price`

## Stability Guarantees

VelesMist will not remove or rename fields inside `velesmist.scan.v1` without introducing a new schema version. New optional fields may be added only when they do not break existing automation.

Unmarketable inventory items are not emitted by this contract yet. Reporting them explicitly with the reserved `not_marketable` recommendation is planned as a separate TDD product PR.
