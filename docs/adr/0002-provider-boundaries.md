# ADR 0002: Provider Boundaries

## Status

Accepted

## Context

The scanner needs to support real Steam inventory reads, deterministic fixtures, and local price data without mixing HTTP, CLI, parsing, scoring, and reporting concerns.

## Decision

Keep scan orchestration in `internal/app` behind provider interfaces:

- inventory providers fetch and parse inventory data;
- price providers load local price data;
- the scorer classifies already-normalized items;
- report code formats a completed scan result.

## Consequences

- CLI code remains an adapter.
- Tests can replace Steam with fixtures or mock HTTP servers.
- Future providers can be added without changing report formatting.
- Provider interfaces should stay narrow and should not expose credentials or session concepts.
