# Repository Ruleset

Target repository: `asketmc/VelesMist`

GitHub rulesets and push protection are repository settings. They cannot be fully enforced from files alone.

## Current Observed Settings

Verified through the GitHub repository API on 2026-07-03:

- secret scanning: enabled;
- secret scanning push protection: enabled;
- private vulnerability reporting: enabled;
- Dependabot security updates: enabled;
- classic `main` branch protection: enabled;
- required pull request reviews: enabled;
- required approving review count: `0` for solo-maintainer mode;
- Code Owners review requirement: disabled for solo-maintainer mode;
- stale review dismissal: enabled;
- required status checks must be up to date: enabled;
- linear history: enabled;
- conversation resolution: enabled;
- force pushes: blocked;
- branch deletions: blocked.

The repository currently has classic branch protection and no repository rulesets.

## Bootstrap Order

1. Push the repository to `main`.
2. Let GitHub Actions run once so check names are selectable.
3. Enable secret scanning and push protection.
4. Create the `main` ruleset or classic branch protection rule.
5. Mark required checks.

## Required `main` Rules

Enable:

- require pull request before merging;
- require approvals only when there is a second human reviewer available;
- require review from Code Owners when reviewer coverage exists;
- dismiss stale approvals when new commits are pushed;
- require status checks to pass before merging;
- require branches to be up to date before merging;
- require conversation resolution before merging;
- block force pushes;
- block branch deletion;
- require linear history if using squash/rebase merge only.

Required checks currently enforced:

- `CI tests`;
- `Analyze Go`;
- `dependency-review`;
- `reuse`;
- `scan`;
- `semgrep`;
- `docs`;
- `sbom`.

Optional after stabilization:

- `Scorecard`.

## Repository Settings

Enable:

- Actions restricted to GitHub Actions and selected pinned third-party actions used here;
- workflow permissions read-only by default;
- Dependabot alerts;
- Dependabot security updates;
- secret scanning;
- push protection;
- code scanning alerts;
- labels: `security`, `good first issue`, `dependencies`, `release`, `documentation`.
