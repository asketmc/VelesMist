# Changelog

## [Unreleased]

## [0.3.1] - 2026-07-03

### Fixed

- Fixed real Steam inventory scans for larger Dota 2 inventories by replacing the rejected `count=5000` request with safe paginated requests.
- Accepted Steam's numeric `more_items` pagination flag in addition to boolean values.
- Clarified HTTP 400 handling as a rejected Steam inventory request instead of implying the inventory is private.

## [0.3.0] - 2026-07-03

### Added

- Minimal localhost web UI via `velesmist ui`.
- UI contract tests for localhost binding, scan API behavior, sanitized errors, and offline fixture scans.
- Focused `make coverage-ui` gate requiring at least 80% statement coverage for `cmd/velesmist`.

### Changed

- `make verify` now includes the focused UI/CLI coverage gate.

## [0.2.0] - 2026-07-03

### Added

- Read-only scan recommendation slice with `sell`, `skip`, and `missing_price` output.
- Local price cache schema `velesmist.price-cache.v1`.
- QA evidence maps, output contracts, ADRs, and local verification gates.

### Fixed

- Reduced elevated GitHub Actions permissions to the jobs that need them.
- Replaced the unpinned REUSE workflow install with a pinned REUSE action.
- Clarified that the README Scorecard badge is workflow status, not a numeric score claim.
- Aligned repository ruleset documentation with the current solo-maintainer branch protection settings.
- Documented GitHub private vulnerability reporting after enabling it in repository settings.

## [0.1.0] - 2026-07-03

### Added

- Standalone Go CLI with `scan` and `version` commands.
- Steam public inventory client with timeout and typed upstream errors.
- Local JSON cache with schema versioning.
- Stable JSON and table reports.
- Unit, integration, golden, fuzz, and assurance tests.
- OSS assurance docs and GitHub Actions supply-chain workflows.

[0.3.1]: https://github.com/asketmc/VelesMist/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/asketmc/VelesMist/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/asketmc/VelesMist/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/asketmc/VelesMist/releases/tag/v0.1.0
