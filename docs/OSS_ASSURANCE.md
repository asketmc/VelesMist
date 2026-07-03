# OSS Assurance

This project implements controls for supply-chain and open-source hygiene. It is not certified secure and has not undergone a formal external audit.

| Control | Status | Evidence |
| --- | --- | --- |
| OpenSSF Scorecard | Partial | `.github/workflows/scorecard.yml`; latest observed score: 6.3/10 on `ac546c856870f1ef9f8bae5eec3ee19131651ef6` |
| OpenSSF Best Practices Badge | Planned | Requires external OpenSSF Best Practices project registration; no placeholder badge is used |
| GitHub Actions CI | Implemented | `.github/workflows/ci.yml` |
| CodeQL | Implemented | `.github/workflows/codeql.yml` |
| Dependabot | Implemented | `.github/dependabot.yml` |
| Dependency Review | Implemented | `.github/workflows/dependency-review.yml` |
| Docs checks | Implemented | `.github/workflows/docs.yml`, `internal/assurance` |
| Format check | Implemented | `make format`, `.github/workflows/ci.yml` |
| Lint / static analysis | Implemented | `make lint`, Staticcheck pinned to `v0.7.0`, `.github/workflows/ci.yml` |
| Typecheck | Implemented | `make typecheck`, `.github/workflows/ci.yml` |
| Coverage artifact | Implemented | `make coverage`, CI uploads `coverage.out` and `coverage.txt` artifacts |
| Coverage badge | Not planned | No external coverage SaaS or generated Pages badge is configured; avoid misleading vanity badges |
| Mutation testing | Not planned | Deferred until stable core APIs justify the extra runtime cost |
| Secret Scanning / Push Protection | Implemented in GitHub setting | Verified enabled through the GitHub repository API on 2026-07-03 |
| REUSE compliance | Implemented | `REUSE.toml`, `LICENSES/MIT.txt`, `.github/workflows/reuse.yml` |
| SPDX license metadata | Implemented | SPDX headers and REUSE annotations |
| SPDX SBOM | Implemented | `.github/workflows/sbom.yml` produces CI artifacts; `v0.1.0` release includes `sbom.spdx.json` |
| CycloneDX SBOM | Implemented | `.github/workflows/sbom.yml` produces CI artifacts; `v0.1.0` release includes `sbom.cdx.json` |
| OSV Scanner | Implemented | `.github/workflows/osv-scanner.yml` |
| govulncheck | Implemented | `Makefile`, `.github/workflows/ci.yml` |
| Semgrep | Implemented | `.github/workflows/semgrep.yml` |
| SLSA provenance / artifact attestations | Implemented | `.github/workflows/release.yml`; `v0.1.0` artifacts verified with `gh attestation verify` |
| Sigstore / cosign release signing | Implemented | `.github/workflows/release.yml`; `v0.1.0` release includes `.sigstore.json` bundles |
| Workflow pinning policy | Implemented | `docs/WORKFLOW_PINNING.md`, `internal/assurance` tests |
| Security Insights metadata | Implemented | `security-insights.yml` uses OpenSSF Security Insights schema `2.2.0`; no fake contact email or unverified release claim |
| Branch protection / required checks | Implemented in GitHub setting | Verified classic `main` branch protection through the GitHub repository API on 2026-07-03 |
| CODEOWNERS review control | Implemented and enforced in GitHub setting | `.github/CODEOWNERS`; branch protection requires Code Owners review |
| Good first issue labels | Implemented in GitHub setting | Verified labels: `good first issue`, `security`, `dependencies`, `release`, `documentation` |
| All Contributors | Implemented | `.all-contributorsrc`, README badge |
| pre-commit.ci | Planned | Requires installing the external pre-commit.ci GitHub app; no placeholder badge is used |
| Docker image / container scan | Not applicable | No container distribution strategy |
| SonarCloud / Codecov / FOSSA / Snyk | Not planned | Avoid external SaaS badges without a concrete need |
| Renovate | Not planned | Dependabot is the selected dependency bot |

## Runtime Dependency Position

VelesMist uses the Go standard library for runtime code. Tooling dependencies are limited to CI/dev workflows and do not ship with the binary.

## Required Checks For `main`

Required checks currently enforced on `main`:

- `CI tests`;
- `Analyze Go`;
- `dependency-review`;
- `reuse`;
- `scan`;
- `semgrep`;
- `docs`;
- `sbom`.

Scorecard should run continuously. Promote it to a required check only after it proves stable for this repository.

## Current Scorecard Notes

The latest observed OpenSSF Scorecard run completed successfully, but the numeric score is not yet a high-score claim. Known non-code or staged gaps:

- OpenSSF Best Practices badge requires external project registration.
- `Maintained` is temporarily low because the repository was created within the last 90 days.
- `Contributors` is low until the project has contributors from more organizations.
- `Code-Review` improves only through future human-reviewed PR history.
- `Branch-Protection` reported an authentication limitation in Scorecard despite branch protection being enabled in GitHub settings.
