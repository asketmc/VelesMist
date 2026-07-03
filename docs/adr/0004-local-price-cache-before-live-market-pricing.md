# ADR 0004: Local Price Cache Before Live Market Pricing

## Status

Accepted

## Context

Live market pricing introduces rate limits, parsing risk, availability risk, and possible terms-of-service ambiguity. The first useful workflow can be served by local/manual price inputs.

## Decision

Use `velesmist.price-cache.v1` as the first price source. The scanner classifies items with local prices and marks marketable items without prices as `missing_price`.

## Consequences

- Tests remain deterministic and network-free.
- Users can audit and reproduce the exact price inputs behind a report.
- Live market pricing remains a future feature and must be introduced through a separate design and test plan.
- Price cache schema changes require contract updates and compatibility review.
