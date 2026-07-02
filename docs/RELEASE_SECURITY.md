# Release Security

## Release Model

Releases are built by `.github/workflows/release.yml` on GitHub release publication.

The workflow:

- builds standalone `CGO_ENABLED=0` binaries for Linux, Windows, and macOS;
- archives binaries with clear OS/architecture names;
- creates `SHA256SUMS.txt`;
- creates SPDX and CycloneDX SBOMs;
- signs artifacts with Sigstore/cosign keyless signing;
- creates GitHub artifact attestations for provenance;
- uploads artifacts to the GitHub release.

## Required Release Gates

Before publishing a release:

- CI passes;
- CodeQL has no unresolved critical findings;
- Dependency Review passes for the release PR;
- REUSE passes;
- govulncheck passes;
- OSV Scanner passes or findings are triaged.

## Not Applicable

Docker image signing and container scanning are not applicable until the project publishes a container image.
