# Artifact Traceability

This document maps generated or expected assurance artifacts to their producer, publication location, verification path, and current status. These are evidence artifacts, not certification claims.

| Artifact | Source workflow or command | Produced when | Published where | Verification | Current status |
| --- | --- | --- | --- | --- | --- |
| Release binary archives | `.github/workflows/release.yml` | GitHub Release is published | GitHub Release assets, for example `v0.1.0` | Download archive, verify checksum, run `velesmist version` | implemented |
| `SHA256SUMS.txt` | `.github/workflows/release.yml` | GitHub Release is published | GitHub Release assets | `sha256sum --check SHA256SUMS.txt` or compare with `Get-FileHash` | implemented |
| SPDX SBOM | `.github/workflows/sbom.yml`, `.github/workflows/release.yml` | PR/push CI and GitHub Release | CI artifacts and GitHub Release asset `sbom.spdx.json` | Parse JSON and inspect package/repository metadata | implemented |
| CycloneDX SBOM | `.github/workflows/sbom.yml`, `.github/workflows/release.yml` | PR/push CI and GitHub Release | CI artifacts and GitHub Release asset `sbom.cdx.json` | Parse JSON and inspect component/dependency metadata | implemented |
| GitHub artifact attestation / provenance | `.github/workflows/release.yml` with `actions/attest-build-provenance` | GitHub Release is published | GitHub artifact attestation service | `gh attestation verify <artifact> --repo asketmc/VelesMist --source-ref refs/tags/<tag>` | implemented |
| Sigstore / cosign bundles | `.github/workflows/release.yml` with `cosign sign-blob` | GitHub Release is published | GitHub Release assets as `*.sigstore.json` | `cosign verify-blob --bundle <asset>.sigstore.json <asset>` | implemented |
| Coverage output | `make coverage` | Developer runs local coverage command | Local `coverage.out` and `coverage.txt`; not published as a SaaS badge | `go tool cover -func=coverage.out` | implemented locally |
| Scorecard result | `.github/workflows/scorecard.yml` | Push to `main`, schedule, branch protection change | GitHub Actions run and OpenSSF published result | Inspect workflow run logs and Scorecard published output | implemented, numeric score not used as a badge |
| CodeQL result | `.github/workflows/codeql.yml` | Push, PR, schedule | GitHub code scanning | Review CodeQL check and code scanning alerts | implemented |
| OSV result | `.github/workflows/osv-scanner.yml` | Push, PR, schedule | GitHub Actions `scan` check | Inspect OSV workflow output | implemented |
| govulncheck result | `.github/workflows/ci.yml`, `make vuln` | PR/push CI or explicit local command | GitHub Actions `CI tests`; local terminal output | CI check or `make vuln` | implemented |
| REUSE result | `.github/workflows/reuse.yml` | Push and PR | GitHub Actions `reuse` check | `reuse lint` in workflow | implemented |
| JSON scan report contract | `schemas/scan-report.v1.json`, `internal/contracts/contracts_test.go` | Tests run | Repository source and CI logs | `go test ./internal/contracts` | implemented |
| Price cache contract | `schemas/price-cache.v1.json`, `internal/contracts/contracts_test.go` | Tests run | Repository source and CI logs | `go test ./internal/contracts` | implemented |

Release verification details are in [VERIFY_RELEASE.md](VERIFY_RELEASE.md).
