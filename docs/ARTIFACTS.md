# Artifact Traceability

This document maps assurance and build artifacts to their producer, trigger, location, verification method, and current status. These artifacts provide control evidence for verification, attestation, and provenance.

| Artifact | Producer | Trigger | Location | Verification method | Status |
| --- | --- | --- | --- | --- | --- |
| Local binary | `make build` / `make verify` | Developer runs local build or verification gate | `dist/velesmist` | Run `dist/velesmist version`; inspect build command for `CGO_ENABLED=0` and `-trimpath` | implemented |
| Release binary archives | `.github/workflows/release.yml` | GitHub Release is published from a tag | GitHub Release assets for the selected tag | Download archive, verify checksum, extract, run `velesmist version` | implemented |
| Checksums | `.github/workflows/release.yml` | GitHub Release is published from a tag | GitHub Release asset `SHA256SUMS.txt` | `sha256sum --check SHA256SUMS.txt` or compare with `Get-FileHash` | implemented |
| SPDX SBOM | `.github/workflows/sbom.yml`, `.github/workflows/release.yml` | PR/push CI and tagged release | CI artifact and GitHub Release asset `sbom.spdx.json` | Parse JSON and inspect package/repository metadata | implemented |
| CycloneDX SBOM | `.github/workflows/sbom.yml`, `.github/workflows/release.yml` | PR/push CI and tagged release | CI artifact and GitHub Release asset `sbom.cdx.json` | Parse JSON and inspect component/dependency metadata | implemented |
| GitHub artifact attestations / provenance | `.github/workflows/release.yml` with `actions/attest-build-provenance` | GitHub Release is published from a tag | GitHub artifact attestation service | `gh attestation verify <artifact> --repo asketmc/VelesMist --source-ref refs/tags/<tag>` | implemented |
| cosign signatures | `.github/workflows/release.yml` with `cosign sign-blob` | GitHub Release is published from a tag | GitHub Release assets as `*.sigstore.json` bundles | `cosign verify-blob --bundle <asset>.sigstore.json <asset>` | implemented |
| `coverage.out` | `make coverage` | Developer runs coverage command | Local workspace root | `go tool cover -func=coverage.out` | implemented |
| Coverage summary | `make coverage` | Developer runs coverage command | `coverage.txt` and terminal output | Inspect package/function summary; no global vanity threshold is enforced | implemented |
| UI/CLI coverage gate | `make coverage-ui` / `make verify` | PR/push CI and local verification | `coverage-ui.out`, `coverage-ui.txt`, and CI logs | Confirm `cmd/velesmist` statement coverage is at least 80% | implemented |
| Scorecard result | `.github/workflows/scorecard.yml` | Push to `main`, schedule, or branch protection change | GitHub Actions run and OpenSSF published result | Inspect workflow logs and Scorecard published output | partial |
| CodeQL result | `.github/workflows/codeql.yml` | Push, PR, or schedule | GitHub code scanning and PR check | Review CodeQL check and code scanning alerts | implemented |
| OSV Scanner result | `.github/workflows/osv-scanner.yml` | Push, PR, or schedule | GitHub Actions `scan` check | Inspect OSV Scanner workflow output | implemented |
| govulncheck result | `.github/workflows/ci.yml`, `make vuln` | PR/push CI or explicit local command | GitHub Actions `CI tests`; local terminal output | CI log or `make vuln`; local run may require network access | implemented |
| Dependency Review result | `.github/workflows/dependency-review.yml` | Pull request | GitHub PR check | Inspect Dependency Review check result | implemented |
| REUSE result | `.github/workflows/reuse.yml` | Push or pull request | GitHub Actions `reuse` check | `reuse lint` in workflow output | implemented |
| Semgrep result | `.github/workflows/semgrep.yml` | Push or pull request | GitHub Actions and code scanning/SARIF output | Inspect Semgrep workflow and code scanning alerts | implemented |
| Docs check result | `.github/workflows/docs.yml` | Push or pull request | GitHub Actions `docs` check | `go test ./internal/assurance` in workflow output | implemented |
| Scan report contract | `schemas/scan-report.v1.json`, `internal/contracts/contracts_test.go` | Tests run | Repository source and CI logs | `go test ./internal/contracts` | implemented |
| Price cache contract | `schemas/price-cache.v1.json`, `internal/contracts/contracts_test.go` | Tests run | Repository source and CI logs | `go test ./internal/contracts` | implemented |

Release verification details are in [VERIFY_RELEASE.md](VERIFY_RELEASE.md).
