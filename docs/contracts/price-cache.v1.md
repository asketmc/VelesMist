# Price Cache Contract v1

Machine-readable schema: [`schemas/price-cache.v1.json`](../../schemas/price-cache.v1.json).

The price cache is a local/manual input format. It does not contain Steam credentials, cookies, session data, or generated sell reports.

## Required Fields

- `schema_version`: must be `velesmist.price-cache.v1`.
- `currency`: currency label used by the report, default `USD`.
- `prices`: object keyed by Steam `market_hash_name`.

## Price Entry Fields

Each `prices` entry must include at least one of:

- `buyer_price_cents`: gross buyer-facing price in cents.
- `lowest_price`: Steam-style money string, for example `$12.34`.
- `median_price`: Steam-style money string, for example `$12.00`.

Optional:

- `source`: human-readable source label such as `manual` or `manual-steam-market-check`.

## Invalid Or Missing Price Behavior

- Unsupported `schema_version` fails with invalid input.
- Unknown top-level fields fail decoding.
- Entries with no usable positive price are ignored by the loader.
- Items without a matching price cache entry appear in scan output as `missing_price`.
