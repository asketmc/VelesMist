# Threat Model

## Assets

- Steam64 identifiers.
- Public inventory data and estimated item values.
- Local cache file.
- Local price cache.
- Release artifacts.
- GitHub Actions workflows and repository settings.

## Trust Boundaries

- Steam Community API responses are external and untrusted.
- Local JSON files are operator-controlled and may be malformed.
- GitHub Actions runners are external CI infrastructure.
- Release consumers trust GitHub release assets only after checksum, signature, and attestation verification.

## Threats And Controls

| Threat | Impact | Controls |
| --- | --- | --- |
| Steam identifier leakage | Privacy loss | Privacy docs; no telemetry; no third-party uploads of private reports |
| API key/token leakage | Account compromise | No API keys/cookies/sessions accepted; issue templates warn against secrets; push protection required |
| Malicious dependency update | Supply-chain compromise | No runtime dependencies; Dependabot; Dependency Review; govulncheck; OSV Scanner; CODEOWNERS |
| Compromised GitHub Actions workflow | Release compromise | CODEOWNERS for `.github/workflows/*`; pinned actions; minimal permissions; Scorecard |
| Tampered release artifact | User compromise | SHA256SUMS; Sigstore keyless signing; GitHub artifact attestations; verification docs |
| Poisoned cache | Misleading report or parse failure | Cache schema version; JSON decoding; cache TTL; no dynamic execution |
| Untrusted inventory/market API response | Crash or misleading report | Typed parser; fuzz test; bounded HTTP response read; no shell execution |
| Steam rate limiting | Incomplete scan | Typed `rate_limited` error mapped to exit code `3`; local cache |

## Security Invariants

- Production code must not require Steam credentials.
- Production code must not execute shell commands.
- Production code must not use CGO.
- Tests must not require network or real Steam credentials.
- Release artifacts must be verifiable from source repository to binary.
