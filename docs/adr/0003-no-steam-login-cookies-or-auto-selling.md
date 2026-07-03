# ADR 0003: No Steam Login, Cookies, Or Auto-Selling

## Status

Accepted

## Context

Steam login, cookies, session scraping, browser automation, market listing creation, and Steam Guard automation would materially increase user risk and project abuse potential.

## Decision

VelesMist is read-only. It does not accept Steam passwords, cookies, session IDs, API keys, or Steam Guard codes. It does not create listings, confirm trades, or automate selling.

## Consequences

- Public/readable inventory scanning is the supported online path.
- Private inventories cannot be scanned through login workarounds.
- Reports are advisory and require manual user action outside VelesMist.
- Issues and test fixtures must not include real credentials, cookies, tokens, private inventory exports, generated sell reports, or local cache files.
