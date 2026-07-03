# VelesMist

[![CI](https://github.com/asketmc/VelesMist/actions/workflows/ci.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/ci.yml)
[![CodeQL](https://github.com/asketmc/VelesMist/actions/workflows/codeql.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/codeql.yml)
[![Dependency Review](https://github.com/asketmc/VelesMist/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/dependency-review.yml)
[![Docs](https://github.com/asketmc/VelesMist/actions/workflows/docs.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/docs.yml)
[![OpenSSF Scorecard](https://github.com/asketmc/VelesMist/actions/workflows/scorecard.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/scorecard.yml)
[![REUSE](https://github.com/asketmc/VelesMist/actions/workflows/reuse.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/reuse.yml)
[![SBOM](https://github.com/asketmc/VelesMist/actions/workflows/sbom.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/sbom.yml)
[![OSV Scanner](https://github.com/asketmc/VelesMist/actions/workflows/osv-scanner.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/osv-scanner.yml)
[![Semgrep](https://github.com/asketmc/VelesMist/actions/workflows/semgrep.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/semgrep.yml)
[![Release Security](https://github.com/asketmc/VelesMist/actions/workflows/release.yml/badge.svg)](https://github.com/asketmc/VelesMist/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

VelesMist is a standalone Go CLI for read-only Steam/Dota 2 inventory analysis. It scans a public Steam inventory, aggregates marketable items, optionally applies a local price cache, and emits stable JSON or human-readable table output for manual review.

It does not log in, does not accept Steam cookies, does not create market listings, and does not automate Steam Guard confirmations.

## Supported Platforms

Release builds target:

- Linux amd64 and arm64;
- Windows amd64 and arm64;
- macOS amd64 and arm64.

Runtime model: one static Go binary built with `CGO_ENABLED=0`. Docker, Python, JVM, Node, and Electron are not required at runtime.

## Install From Releases

Download the archive for your OS and architecture from:

```text
https://github.com/asketmc/VelesMist/releases
```

Before running a release binary, verify checksums, Sigstore signing, GitHub artifact attestation, and SBOM artifacts with [docs/VERIFY_RELEASE.md](docs/VERIFY_RELEASE.md).

## Build From Source

```bash
git clone https://github.com/asketmc/VelesMist.git
cd VelesMist
make test
make build
```

Requires Go 1.26 or newer for source builds.

Equivalent build command:

```bash
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o dist/velesmist ./cmd/velesmist
```

Release builds inject version metadata with ldflags:

- `github.com/asketmc/VelesMist/internal/version.Version`;
- `github.com/asketmc/VelesMist/internal/version.Commit`;
- `github.com/asketmc/VelesMist/internal/version.BuildDate`;
- `github.com/asketmc/VelesMist/internal/version.Dirty`.

## Usage

```bash
velesmist scan --steam-id 76561198000000000
velesmist scan --steam-id 76561198000000000 --format table
velesmist scan --steam-id 76561198000000000 --format json
velesmist version
```

Optional price cache:

```bash
velesmist scan \
  --steam-id 76561198000000000 \
  --format json \
  --price-cache prices.json \
  --min-price 5.00
```

Price cache format:

```json
{
  "Golden Moonfall": {
    "lowest_price": "$12.34"
  },
  "Jagged Honor | Blade": {
    "buyer_price_cents": 400
  }
}
```

## JSON Output

JSON output is intended to be machine-readable and stable within the `velesmist.scan.v1` schema:

```json
{
  "schema_version": "velesmist.scan.v1",
  "steam_id": "76561198000000000",
  "appid": 570,
  "contextid": "2",
  "currency": "USD",
  "threshold_cents": 500,
  "items": [],
  "candidates": [],
  "summary": {
    "marketable_items": 0,
    "priced_items": 0,
    "candidate_items": 0,
    "estimated_total_receive_cents": 0
  }
}
```

## Exit Codes

| Code | Meaning |
| ---: | --- |
| 0 | Success |
| 1 | Internal/runtime error |
| 2 | Invalid input or configuration |
| 3 | Upstream Steam API unavailable, timed out, or rate-limited |

## Cache

Inventory HTTP responses are cached locally as JSON.

- Schema version: `1`.
- Default location: OS user cache directory under `velesmist/cache.json`.
- Override: `--cache-file <path>`.
- Disable: `--no-cache`.
- TTL: `--cache-ttl`, default `10m`.

Delete the cache file to clear local state.

## Privacy And Security

VelesMist reads a Steam64 ID, public inventory metadata, optional local price cache, and local cache file. It sends read-only HTTPS requests to Steam Community public inventory endpoints. It does not require API keys, Steam credentials, cookies, session IDs, or payment data.

See:

- [SECURITY.md](SECURITY.md)
- [docs/PRIVACY.md](docs/PRIVACY.md)
- [docs/THREAT_MODEL.md](docs/THREAT_MODEL.md)
- [docs/OSS_ASSURANCE.md](docs/OSS_ASSURANCE.md)

## Development

```bash
make test
make lint
make vet
make vuln
make build
make snapshot-release
```

Tests use fixtures and mock HTTP servers. They do not require network access or Steam credentials.

## License

MIT. See [LICENSE](LICENSE).
