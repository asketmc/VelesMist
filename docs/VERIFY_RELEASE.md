# Verify a Release

Use these steps to verify release artifacts from `asketmc/VelesMist`.

Replace `v0.1.0` with the release tag.

## Download Artifacts

Download from:

```text
https://github.com/asketmc/VelesMist/releases
```

Expected files include:

- `velesmist-v0.1.0-linux-amd64.tar.gz`;
- `velesmist-v0.1.0-linux-arm64.tar.gz`;
- `velesmist-v0.1.0-windows-amd64.zip`;
- `velesmist-v0.1.0-windows-arm64.zip`;
- `velesmist-v0.1.0-darwin-amd64.tar.gz`;
- `velesmist-v0.1.0-darwin-arm64.tar.gz`;
- `SHA256SUMS.txt`;
- `*.sigstore.json` bundles;
- `sbom.spdx.json`;
- `sbom.cdx.json`.

## Verify Checksum

Linux/macOS:

```bash
sha256sum --check SHA256SUMS.txt
```

PowerShell:

```powershell
Get-FileHash .\velesmist-v0.1.0-windows-amd64.zip -Algorithm SHA256
Get-Content .\SHA256SUMS.txt
```

## Verify Sigstore Signature

Example for the checksum file:

```bash
cosign verify-blob \
  --bundle SHA256SUMS.txt.sigstore.json \
  --certificate-identity-regexp "https://github.com/asketmc/VelesMist/.github/workflows/release.yml@refs/tags/.+" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  SHA256SUMS.txt
```

Repeat for individual archives when desired.

## Verify GitHub Artifact Attestation

```bash
gh attestation verify velesmist-v0.1.0-linux-amd64.tar.gz --repo asketmc/VelesMist
```

This verifies the provenance connection between the artifact, workflow, and source repository.

## Inspect SBOMs

Confirm both files are present:

- `sbom.spdx.json`;
- `sbom.cdx.json`.

At minimum, confirm package name, version, repository URL, and dependency set match the release notes.

## Failure Policy

Do not install or redistribute a release if:

- checksum verification fails;
- Sigstore verification fails;
- GitHub attestation verification fails;
- SBOM files are missing from a release that claims to provide them;
- release assets include real inventory exports, generated sell reports, cookies, tokens, API keys, or cache files.
