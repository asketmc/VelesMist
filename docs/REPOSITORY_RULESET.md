# Repository Ruleset

Target repository: `asketmc/VelesMist`

GitHub rulesets and push protection are repository settings. They cannot be fully enforced from files alone.

## Bootstrap Order

1. Push the repository to `main`.
2. Let GitHub Actions run once so check names are selectable.
3. Enable secret scanning and push protection.
4. Create the `main` ruleset or classic branch protection rule.
5. Mark required checks.

## Required `main` Rules

Enable:

- require pull request before merging;
- require at least one approval;
- require review from Code Owners;
- dismiss stale approvals when new commits are pushed;
- require status checks to pass before merging;
- require branches to be up to date before merging;
- require conversation resolution before merging;
- block force pushes;
- block branch deletion;
- require linear history if using squash/rebase merge only.

Required checks:

- `CI tests`;
- `CodeQL`;
- `Dependency Review`;
- `Docs`;
- `REUSE`;
- `OSV Scanner`.

Optional after stabilization:

- `Semgrep`;
- `Scorecard`;
- `SBOM`.

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
