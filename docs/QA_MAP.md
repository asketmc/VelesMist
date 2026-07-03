# QA Map

This map ties product requirements and risks to concrete test evidence. Status values are limited to `implemented`, `partial`, `missing`, `planned`, and `not applicable`.

| Area | Requirement | Risk | Evidence / test file | Validation command | Status |
| --- | --- | --- | --- | --- | --- |
| CLI | CLI flag validation | Bad input silently produces misleading reports or wrong exit codes | `internal/config/config_test.go`, `cmd/velesmist/main_test.go` | `go test ./internal/config ./cmd/velesmist` | implemented |
| CLI | `--fixture` offline scan path | Offline demos/tests accidentally call Steam or require credentials | `cmd/velesmist/main_test.go`, `internal/inventory/fixture_provider.go` | `go test ./cmd/velesmist` | implemented |
| CLI | `--game dota2` handling | Unsupported game/app/context combinations are accepted silently | `internal/config/config_test.go`, `cmd/velesmist/main_test.go` | `go test ./internal/config ./cmd/velesmist` | implemented |
| Pricing | `price-cache` `schema_version` | Old or malformed local price files are accepted silently | `internal/pricing/pricing_test.go`, `internal/contracts/contracts_test.go`, `schemas/price-cache.v1.json` | `go test ./internal/pricing ./internal/contracts` | implemented |
| Pricing | Price cache parsing | Manual prices are misread or schema drift changes scoring | `internal/pricing/pricing_test.go`, `internal/contracts/contracts_test.go` | `go test ./internal/pricing ./internal/contracts` | implemented |
| Pricing | Missing price behavior | Items with no local price are incorrectly sold or skipped | `internal/pricing/pricing_test.go`, `internal/report/testdata/scan.json.golden` | `go test ./internal/pricing ./internal/report ./internal/contracts` | implemented |
| Pricing | `sell` / `skip` classification | Items cross the threshold in the wrong direction | `internal/pricing/pricing_test.go`, `cmd/velesmist/main_test.go` | `go test ./internal/pricing ./cmd/velesmist` | implemented |
| Pricing | Fee calculation | Estimated platform fee is wrong | `internal/pricing/pricing_test.go` | `go test ./internal/pricing` | implemented |
| Pricing | Seller receive calculation | Seller proceeds are wrong and the sell threshold becomes unreliable | `internal/pricing/pricing_test.go` | `go test ./internal/pricing` | implemented |
| Report | `reason_codes` | Reports become hard to audit or explain | `internal/pricing/pricing_test.go`, `internal/report/report_test.go`, `internal/contracts/contracts_test.go` | `go test ./internal/pricing ./internal/report ./internal/contracts` | implemented |
| Report | `confidence` | Automation cannot distinguish local-cache prices from missing prices | `internal/pricing/pricing_test.go`, `internal/contracts/contracts_test.go` | `go test ./internal/pricing ./internal/contracts` | implemented |
| Report | `liquidity_score` | Consumers infer live liquidity where none exists | `internal/report/testdata/scan.json.golden`, `schemas/scan-report.v1.json` | `go test ./internal/report ./internal/contracts` | implemented |
| Report | `market_url` | Users cannot verify item pages manually | `internal/pricing/pricing_test.go`, `internal/report/testdata/scan.json.golden` | `go test ./internal/pricing ./internal/report ./internal/contracts` | implemented |
| Report | JSON output stability | Automation breaks due to unreviewed output changes | `internal/report/report_test.go`, `internal/report/testdata/scan.json.golden`, `internal/contracts/contracts_test.go`, `schemas/scan-report.v1.json` | `go test ./internal/report ./internal/contracts` | implemented |
| Report | Table output stability | Human-readable review output changes unexpectedly | `internal/report/report_test.go`, `internal/report/testdata/scan.table.golden` | `go test ./internal/report` | implemented |
| Report | Sorting order | High-value candidates are hidden below less actionable rows | `internal/pricing/pricing_test.go`, `internal/report/testdata/scan.json.golden` | `go test ./internal/pricing ./internal/report` | implemented |
| Errors | Exit code mapping | Automation cannot distinguish invalid input from upstream failures | `internal/errors/errors_test.go`, `cmd/velesmist/main_test.go` | `go test ./internal/errors ./cmd/velesmist` | implemented |
| Steam provider | Mocked Steam provider behavior | Unit tests become flaky, leak identifiers, or depend on Steam availability | `internal/steam/client_test.go`, `internal/steam/provider_test.go` | `go test ./internal/steam` | implemented |
| Test isolation | No real network in tests | Tests depend on real Steam or network availability | `internal/steam/client_test.go`, `internal/steam/provider_test.go`, `cmd/velesmist/main_test.go` | `go test ./...` | implemented |
| Test isolation | No credentials in tests | Fixtures or logs leak Steam cookies, API keys, session data, or private exports | `cmd/velesmist/main_test.go`, `internal/steam/client_test.go`, `internal/inventory/testdata/dota_inventory.json` | `go test ./...` plus PR checklist review | implemented |
| QA evidence | Coverage command without vanity thresholding | A global percentage target encourages meaningless tests instead of risk coverage | `Makefile`, `docs/ARTIFACTS.md`, this requirement map | `make coverage` | implemented |
| Build | CGO-free build | Release/runtime model gains hidden native dependencies | `Makefile` | `make verify`; `CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o dist/velesmist ./cmd/velesmist` | implemented |
| Documentation | Privacy/security docs coverage | Users misunderstand Steam ID, cache, credential, or reporting behavior | `README.md`, `SECURITY.md`, `docs/PRIVACY.md`, `docs/THREAT_MODEL.md`, `.github/pull_request_template.md` | `go test ./internal/assurance`; `make verify` | partial |
| Future behavior | Unmarketable item reporting | Inventory items that cannot be sold disappear from output | Not covered yet | Future TDD PR | planned |
| Future behavior | Live market pricing | Users may expect current prices instead of local/manual cache prices | Not covered yet | Future design/test PR | planned |

`make verify` is the normal deterministic local gate. Network/security checks such as `make vuln` stay separate because they may need module or vulnerability database access.
