# QA Map

This map ties VelesMist requirements and product risks to concrete tests or verification commands. Status values are `implemented`, `partial`, `missing`, or `planned`.

| Requirement | Risk | Test file(s) | Command | Status |
| --- | --- | --- | --- | --- |
| CLI flag validation | Bad input silently produces misleading reports or wrong exit codes | `internal/config/config_test.go`, `cmd/velesmist/main_test.go` | `go test ./internal/config ./cmd/velesmist` | implemented |
| Fixture scan path | Offline demos/tests accidentally call Steam or require credentials | `cmd/velesmist/main_test.go`, `internal/inventory/fixture_provider.go` | `go test ./cmd/velesmist` | implemented |
| Steam provider behavior via mocked HTTP only | Unit tests become flaky, leak identifiers, or depend on Steam availability | `internal/steam/client_test.go`, `internal/steam/provider_test.go` | `go test ./internal/steam` | implemented |
| Price cache parsing | Manual prices are misread or schema drift changes scoring | `internal/pricing/pricing_test.go`, `internal/contracts/contracts_test.go`, `schemas/price-cache.v1.json` | `go test ./internal/pricing ./internal/contracts` | implemented |
| Price cache schema versioning | Old or malformed local price files are accepted silently | `internal/pricing/pricing_test.go`, `internal/contracts/contracts_test.go` | `go test ./internal/pricing ./internal/contracts` | implemented |
| Fee calculation | Gross buyer price, estimated fee, or seller receive are wrong | `internal/pricing/pricing_test.go` | `go test ./internal/pricing` | implemented |
| Recommendation classification | Items are placed into the wrong `sell`, `skip`, or `missing_price` bucket | `internal/pricing/pricing_test.go`, `cmd/velesmist/main_test.go` | `go test ./internal/pricing ./cmd/velesmist` | implemented |
| `reason_codes` | Reports become hard to audit or explain | `internal/report/report_test.go`, `internal/pricing/pricing_test.go` | `go test ./internal/report ./internal/pricing` | implemented |
| JSON output stability | Automation breaks due to unreviewed output changes | `internal/report/report_test.go`, `internal/report/testdata/scan.json.golden`, `internal/contracts/contracts_test.go`, `schemas/scan-report.v1.json` | `go test ./internal/report ./internal/contracts` | implemented |
| Table output stability | Human review output changes unexpectedly | `internal/report/report_test.go`, `internal/report/testdata/scan.table.golden` | `go test ./internal/report` | implemented |
| Exit code mapping | Automation cannot distinguish invalid input from upstream failures | `internal/errors/errors_test.go`, `cmd/velesmist/main_test.go` | `go test ./internal/errors ./cmd/velesmist` | implemented |
| No credentials/secrets in tests | Fixtures or logs leak Steam cookies, API keys, generated reports, or cache files | `cmd/velesmist/main_test.go`, `internal/steam/client_test.go`, `internal/inventory/testdata/dota_inventory.json` | `go test ./...` plus review checklist | implemented |
| No real network calls in unit tests | Local and CI tests become nondeterministic | `internal/steam/client_test.go`, `internal/steam/provider_test.go`, `cmd/velesmist/main_test.go` | `go test ./...` | implemented |
| CGO-free build | Release/runtime model gains hidden native dependencies | `Makefile` | `CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o dist/velesmist ./cmd/velesmist`; `make verify` | implemented |
| Unmarketable item reporting | Inventory items that cannot be sold disappear from output | Not covered yet | Future TDD PR | planned |
| Live market pricing | Users may expect current prices instead of local/manual cache prices | Not covered yet | Future design/test PR | planned |

`make verify` is the local non-network quality gate. Network/security scans remain explicit commands such as `make vuln` and the GitHub OSV workflow.
