<!-- SPDX-FileCopyrightText: 2026 VelesMist contributors -->
<!-- SPDX-License-Identifier: MIT -->

## Summary

- 

## Product Behavior

- [ ] Product behavior changed
- [ ] Product behavior unchanged
- If changed, describe the user-visible behavior:

## Output Contract

- [ ] JSON/table output contract changed
- [ ] JSON/table output contract unchanged
- If changed, update `schemas/*`, golden tests, and contract docs.

## Security / Privacy Impact

- [ ] No new external network path
- [ ] No new secret, token, cookie, session, or credential handling
- [ ] Privacy docs updated when data flow changes
- [ ] Security-sensitive files have CODEOWNERS coverage

## Validation

- [ ] `make verify`
- [ ] `make test`
- [ ] `make lint`
- [ ] `make vet`
- [ ] `make coverage`
- [ ] `make vuln` if dependency/security risk changed
- [ ] Docs updated for user-visible behavior changes
- [ ] GitHub checks are green

## Safety

- [ ] No Steam credentials, cookies, session IDs, API keys, private inventory exports, generated sell reports, or local cache files are included
- [ ] No auto-listing, auto-selling, or Steam Guard automation is introduced
- [ ] No browser automation, scraping, prediction/ML, or Docker runtime is introduced

## AI Disclosure

```text
AI-Assisted: yes/no
Tool: <tool/model name if applicable>
Scope: <tests/docs/code/refactor>
Human-reviewed: yes/no
```
