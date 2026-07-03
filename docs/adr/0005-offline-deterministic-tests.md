# ADR 0005: Offline Deterministic Tests

## Status

Accepted

## Context

Inventory APIs can be unavailable, rate-limited, private, or change response shape. Tests that depend on the live Steam service would be slow, flaky, and risky for privacy.

## Decision

Unit and integration tests must use local fixtures or mock HTTP servers. They must not require network access, Steam credentials, cookies, API keys, tokens, or private inventory data.

## Consequences

- `--fixture` is part of the QA surface for deterministic scan output.
- Mock HTTP server tests cover Steam client/provider behavior.
- Live Steam verification, if ever needed, must be explicit and outside the default test gate.
- `make verify` remains a local non-network gate.
