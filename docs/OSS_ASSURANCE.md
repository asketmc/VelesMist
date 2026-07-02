# OSS Assurance

This project implements controls for supply-chain and open-source hygiene. It is not certified secure and has not undergone a formal external audit.

| Control | Status | Evidence |
| --- | --- | --- |
| OpenSSF Scorecard | Implemented | `.github/workflows/scorecard.yml` |
| GitHub Actions CI | Implemented | `.github/workflows/ci.yml` |
| CodeQL | Implemented | `.github/workflows/codeql.yml` |
| Dependabot | Implemented | `.github/dependabot.yml` |
| Dependency Review | Implemented | `.github/workflows/dependency-review.yml` |
| Docs checks | Implemented | `.github/workflows/docs.yml`, `internal/assurance` |
| Secret Scanning / Push Protection | Implemented | Enabled in GitHub repository settings |
| REUSE compliance | Implemented | `REUSE.toml`, `LICENSES/MIT.txt`, `.github/workflows/reuse.yml` |
| SPDX license metadata | Implemented | SPDX headers and REUSE annotations |
| SPDX SBOM | Implemented | `.github/workflows/sbom.yml`, `.github/workflows/release.yml` |
| CycloneDX SBOM | Implemented | `.github/workflows/sbom.yml`, `.github/workflows/release.yml` |
| OSV Scanner | Implemented | `.github/workflows/osv-scanner.yml` |
| govulncheck | Implemented | `Makefile`, `.github/workflows/ci.yml` |
| Semgrep | Implemented | `.github/workflows/semgrep.yml` |
| SLSA provenance / artifact attestations | Requires first release | `.github/workflows/release.yml` |
| Sigstore / cosign release signing | Requires first release | `.github/workflows/release.yml` |
| Workflow pinning policy | Implemented | `docs/WORKFLOW_PINNING.md`, `internal/assurance` tests |
| Security Insights metadata | Implemented | `security-insights.yml` |
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
