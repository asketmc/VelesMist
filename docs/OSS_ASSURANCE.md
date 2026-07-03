# OSS Assurance

This project implements controls for supply-chain and open-source hygiene. It is not certified secure and has not undergone a formal external audit.

Status values are limited to `implemented`, `partial`, `planned`, `not applicable`, `requires GitHub setting`, and `requires first release`.

| Control | Status | Evidence |
| --- | --- | --- |
| OpenSSF Scorecard | partial | `.github/workflows/scorecard.yml`; the README badge is workflow status, not a numeric score claim |
| CI | implemented | `.github/workflows/ci.yml`; `make verify`, govulncheck, and version smoke run in CI |
| CodeQL | implemented | `.github/workflows/codeql.yml`; Go analysis runs on PR/push/schedule |
| Dependabot | implemented | `.github/dependabot.yml` for GitHub Actions and Go modules |
| Dependency Review | implemented | `.github/workflows/dependency-review.yml`; PR-only dependency/license review |
| Secret scanning / push protection | implemented | GitHub repository settings report secret scanning and push protection enabled |
| REUSE | implemented | `REUSE.toml`, `LICENSES/MIT.txt`, SPDX headers/annotations, `.github/workflows/reuse.yml` |
| SPDX SBOM | implemented | `.github/workflows/sbom.yml`, `.github/workflows/release.yml`, `docs/ARTIFACTS.md` |
| CycloneDX SBOM | implemented | `.github/workflows/sbom.yml`, `.github/workflows/release.yml`, `docs/ARTIFACTS.md` |
| SLSA / GitHub artifact attestations | implemented | `.github/workflows/release.yml`; `v0.1.0` artifacts verified with `gh attestation verify` |
| Sigstore / cosign | implemented | `.github/workflows/release.yml`; release assets include `.sigstore.json` bundles |
| govulncheck | implemented | `.github/workflows/ci.yml`, `make vuln`; kept separate from `make verify` because local runs may need network access |
| OSV Scanner | implemented | `.github/workflows/osv-scanner.yml` |
| Semgrep | implemented | `.github/workflows/semgrep.yml` |
| workflow pinning | implemented | Workflows use commit SHA pins; `docs/WORKFLOW_PINNING.md`; `internal/assurance` tests check pinning and documentation sync |
| Security Insights | implemented | `security-insights.yml` uses OpenSSF Security Insights schema `2.2.0` and avoids placeholder contact data |
| CODEOWNERS | implemented | `.github/CODEOWNERS` documents ownership for sensitive paths; required CODEOWNERS approval is not enforced in solo-maintainer mode |
| branch protection | implemented | `main` requires PR-based flow with required checks, conversation resolution, linear history, blocks force-push, blocks deletion, and enforces admins; approval count is `0` for solo-maintainer mode |
| QA_MAP | implemented | `docs/QA_MAP.md` maps requirements and risks to test evidence |
| ARTIFACTS | implemented | `docs/ARTIFACTS.md` maps artifacts to producers and verification methods |
| output schemas/contracts | implemented | `schemas/scan-report.v1.json`, `schemas/price-cache.v1.json`, `docs/contracts/*`, `internal/contracts` tests |
| release verification docs | implemented | `docs/VERIFY_RELEASE.md` documents checksum, cosign bundle, attestation/provenance, and SBOM verification |
| Docker image / container scan | not applicable | No container distribution strategy |
| SonarCloud / Codecov / Coveralls / FOSSA / Snyk | not applicable | External SaaS badges are intentionally not used for this project slice |
| Renovate | not applicable | Dependabot is the selected dependency bot |

## Runtime Dependency Position

VelesMist uses the Go standard library for runtime code. Tooling dependencies are limited to CI/dev workflows and do not ship with the binary.

## Required Checks For `main`

Required checks on `main`:

- `CI tests`;
- `Analyze Go`;
- `dependency-review`;
- `reuse`;
- `scan`;
- `semgrep`;
- `docs`;
- `sbom`.

Scorecard runs on push, schedule, and branch protection changes. It is not required on pull requests because the workflow is not configured to run on `pull_request`.
