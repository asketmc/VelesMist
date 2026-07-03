# ADR 0006: Contract-First Output

## Status

Accepted.

## Context

VelesMist JSON output is intended for automation as well as human review. Unreviewed field changes, recommendation changes, or price-cache schema drift can break downstream scripts even when the CLI still appears to work.

## Decision

Define stable machine-readable contracts for scan reports and price-cache files before expanding output behavior. Golden output tests and contract tests must fail when JSON fields, recommendation values, or required price-cache fields drift unexpectedly.

## Consequences

Output changes require deliberate schema, fixture, test, and documentation updates. This adds some maintenance cost, but it makes review and release verification more explicit and keeps future product PRs from changing automation contracts accidentally.
