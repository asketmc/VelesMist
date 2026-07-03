# Privacy

VelesMist is a local read-only CLI.

## Data Read

- Steam64 ID supplied with `--steam-id`.
- Public Steam inventory metadata returned by Steam Community.
- Optional local price cache supplied with `--price-cache`.
- Optional local Steam inventory fixture supplied with `--fixture`.
- Local inventory cache file when caching is enabled.

## Data Saved Locally

Inventory HTTP responses may be cached as JSON.

- Schema version: `1`.
- Default location: OS user cache directory under `velesmist/cache.json`.
- Override: `--cache-file <path>`.
- Disable: `--no-cache`.
- TTL: `--cache-ttl`.

Delete the cache file to remove local cached data.

The optional `velesmist.price-cache.v1` price cache is supplied by the operator and stays local. It is read from the path passed with `--price-cache`; VelesMist does not upload it or use it to call a pricing provider.

## Data Sent To External APIs

VelesMist sends read-only HTTPS requests to Steam Community public inventory endpoints.

It sends:

- Steam64 ID;
- appid;
- contextid;
- language and count query parameters.

## Data Not Logged Or Stored

VelesMist does not need and must not log or store:

- Steam password;
- Steam Guard code;
- Steam cookies;
- session IDs;
- API keys;
- payment card data;
- private generated sell reports unless the operator writes output to a file.

## Third-Party Scanning

Only public source release archives are pre-approved for public multi-vendor scanning. Do not upload private inventory exports, local cache files, or generated reports to third-party scanners without explicit owner approval.
