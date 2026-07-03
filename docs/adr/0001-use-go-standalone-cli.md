# ADR 0001: Use Go Standalone CLI

## Status

Accepted

## Context

VelesMist should run as one executable on Windows, Linux, and macOS without requiring Python, JVM, Node, Electron, Docker, or a local service runtime.

## Decision

Implement VelesMist as a Go CLI using Go modules and build release binaries with `CGO_ENABLED=0`.

## Consequences

- Release artifacts can be distributed as simple archives.
- Runtime dependency review is simpler because the binary has no external runtime.
- UI richness is intentionally out of scope for the first versions.
- Cross-platform behavior must be tested through command output, exit codes, and release smoke tests.
