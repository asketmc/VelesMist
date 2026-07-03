# Verify a Release

Use these steps to verify release artifacts from `asketmc/VelesMist`.

Set `TAG` to the release tag you are verifying.

## Download Artifacts

Download from:

```text
https://github.com/asketmc/VelesMist/releases
```

Expected files include:

- `velesmist-${TAG}-linux-amd64.tar.gz`;
- `velesmist-${TAG}-linux-arm64.tar.gz`;
- `velesmist-${TAG}-windows-amd64.zip`;
- `velesmist-${TAG}-windows-arm64.zip`;
- `velesmist-${TAG}-darwin-amd64.tar.gz`;
- `velesmist-${TAG}-darwin-arm64.tar.gz`;
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
$TAG = "v0.2.0"
Get-FileHash ".\velesmist-$TAG-windows-amd64.zip" -Algorithm SHA256
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
TAG=v0.2.0
gh attestation verify "velesmist-${TAG}-linux-amd64.tar.gz" --repo asketmc/VelesMist
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
