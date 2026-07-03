# OSS Assurance

This project implements controls for supply-chain and open-source hygiene. It is not certified secure and has not undergone a formal external audit.

| Control | Status | Evidence |
| --- | --- | --- |
| OpenSSF Scorecard | Partial | `.github/workflows/scorecard.yml`; latest observed score should be read from published Scorecard output, not inferred from the workflow badge |
| GitHub Actions CI | Implemented | `.github/workflows/ci.yml` |
| CodeQL | Implemented | `.github/workflows/codeql.yml` |
| Dependabot | Implemented | `.github/dependabot.yml` |
| Dependency Review | Implemented | `.github/workflows/dependency-review.yml` |
| Docs checks | Implemented | `.github/workflows/docs.yml`, `internal/assurance` |
| Local verification gate | Implemented | `make verify` runs format, tests, vet, and CGO-free build |
| Coverage output | Implemented | `make coverage` produces `coverage.out` and `coverage.txt` locally |
| Output contracts | Implemented | `schemas/scan-report.v1.json`, `schemas/price-cache.v1.json`, `internal/contracts` tests |
| QA evidence map | Implemented | `docs/QA_MAP.md` |
| Artifact traceability | Implemented | `docs/ARTIFACTS.md`, `docs/VERIFY_RELEASE.md` |
| Secret Scanning / Push Protection | Implemented | Enabled in GitHub repository settings |
| REUSE compliance | Implemented | `REUSE.toml`, `LICENSES/MIT.txt`, `.github/workflows/reuse.yml` |
| SPDX license metadata | Implemented | SPDX headers and REUSE annotations |
| SPDX SBOM | Implemented | `.github/workflows/sbom.yml`, `.github/workflows/release.yml` |
| CycloneDX SBOM | Implemented | `.github/workflows/sbom.yml`, `.github/workflows/release.yml` |
| OSV Scanner | Implemented | `.github/workflows/osv-scanner.yml` |
| govulncheck | Implemented | `Makefile`, `.github/workflows/ci.yml` |
| Semgrep | Implemented | `.github/workflows/semgrep.yml` |
| SLSA provenance / artifact attestations | Implemented | `.github/workflows/release.yml`; `v0.1.0` artifacts verified with `gh attestation verify` |
| Sigstore / cosign release signing | Implemented | `.github/workflows/release.yml`; `v0.1.0` release includes `.sigstore.json` bundles |
| Workflow pinning policy | Implemented | `docs/WORKFLOW_PINNING.md`, `internal/assurance` tests |
| Security Insights metadata | Implemented | `security-insights.yml` uses OpenSSF Security Insights schema `2.2.0`; no fake contact email or unverified release claim |
| Branch protection / required checks | Implemented | `main` branch protection and `docs/REPOSITORY_RULESET.md` |
| CODEOWNERS review control | Implemented | `.github/CODEOWNERS`, enforced by `main` branch protection |
| Docker image / container scan | Not applicable | No container distribution strategy |
| SonarCloud / Codecov / FOSSA / Snyk | Not planned | Avoid external SaaS badges without a concrete need |
| Renovate | Not planned | Dependabot is the selected dependency bot |

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
