# Security Policy

## Supported Versions

| Version | Supported |
| --- | --- |
| Latest release | Yes |
| Older releases | Security fixes only when explicitly marked supported |
| Unreleased local drafts | No public security support |

## Reporting a Vulnerability

Use GitHub private vulnerability reporting when the issue is exploitable or may expose sensitive data:

https://github.com/asketmc/VelesMist/security/advisories/new

Do not publish real Steam passwords, Steam Guard codes, Steam cookies, session IDs, API keys, private inventory exports, generated sell reports, or local cache files in public issues.

Include:

- affected version or commit;
- reproduction steps using sanitized fixtures where possible;
- expected impact;
- whether any private Steam identifier, inventory data, or local cache data is involved.

## Supported Use

Supported workflows are read-only public inventory analysis and local report generation.

## Not Supported

The project will not accept features that:

- collect Steam credentials;
- store Steam cookies or sessions;
- automate market listings;
- automate Steam Guard confirmations;
- bypass Steam rate limits;
- add hidden downloads, dynamic code execution, bundled executables, or production subprocess calls.

## Response Target

- Initial maintainer response: 7 days.
- Triage target: 14 days.
- Public advisory or fix notes: after a fix or mitigation exists.
