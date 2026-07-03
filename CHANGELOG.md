# Changelog

## [Unreleased]

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

[0.2.0]: https://github.com/asketmc/VelesMist/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/asketmc/VelesMist/releases/tag/v0.1.0
