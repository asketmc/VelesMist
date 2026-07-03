# Workflow Pinning

GitHub Actions are pinned to commit SHAs. Human-readable tag comments are kept for review.

| Action | Tag | SHA |
| --- | --- | --- |
| `actions/checkout` | `v7.0.0` | `9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0` |
| `actions/setup-go` | `v6.5.0` | `924ae3a1cded613372ab5595356fb5720e22ba16` |
| `actions/upload-artifact` | `v7.0.1` | `043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` |
| `actions/dependency-review-action` | `v5.0.0` | `a1d282b36b6f3519aa1f3fc636f609c47dddb294` |
| `actions/attest-build-provenance` | `v3` | `977bb373ede98d70efdf65b84cb5f73e068dcc2a` |
| `github/codeql-action/init` | `v4` | `54f647b7e1bb85c95cddabcd46b0c578ec92bc1a` |
| `github/codeql-action/autobuild` | `v4` | `54f647b7e1bb85c95cddabcd46b0c578ec92bc1a` |
| `github/codeql-action/analyze` | `v4` | `54f647b7e1bb85c95cddabcd46b0c578ec92bc1a` |
| `github/codeql-action/upload-sarif` | `v4` | `54f647b7e1bb85c95cddabcd46b0c578ec92bc1a` |
| `google/osv-scanner-action/osv-scanner-action` | `v2.3.8` | `9a498708959aeaef5ef730655706c5a1df1edbc2` |
| `ossf/scorecard-action` | `v2.4.3` | `4eaacf0543bb3f2c246792bd56e8cdeffafb205a` |
| `anchore/sbom-action` | `v0` | `e22c389904149dbc22b58101806040fa8d37a610` |
| `CycloneDX/gh-gomod-generate-sbom` | `v2.0.0` | `efc74245d6802c8cefd925620515442756c70d8f` |
| `sigstore/cosign-installer` | `v4.1.2` | `6f9f17788090df1f26f669e9d70d6ae9567deba6` |
| `fsfe/reuse-action` | `v6.0.0` | `676e2d560c9a403aa252096d99fcab3e1132b0f5` |

## Maintenance Rule

When updating an action:

1. Resolve the tag to a commit SHA.
2. Use the SHA in workflow `uses:`.
3. Keep the tag in a comment.
4. Update this table.
5. Run `go test ./...` so `internal/assurance` checks the policy.
