# Workflow Pinning

GitHub Actions are pinned to commit SHAs. Human-readable tag comments are kept for review.

| Action | Tag | SHA |
| --- | --- | --- |
| `actions/checkout` | `v6.0.3` | `df4cb1c069e1874edd31b4311f1884172cec0e10` |
| `actions/setup-go` | `v6.0.0` | `44694675825211faa026b3c33043df3e48a5fa00` |
| `actions/upload-artifact` | `v4` | `ea165f8d65b6e75b540449e92b4886f43607fa02` |
| `actions/dependency-review-action` | `v5.0.0` | `a1d282b36b6f3519aa1f3fc636f609c47dddb294` |
| `actions/attest-build-provenance` | `v3` | `977bb373ede98d70efdf65b84cb5f73e068dcc2a` |
| `github/codeql-action` | `v4` | `54f647b7e1bb85c95cddabcd46b0c578ec92bc1a` |
| `google/osv-scanner-action` | `v2.3.8` | `9a498708959aeaef5ef730655706c5a1df1edbc2` |
| `ossf/scorecard-action` | `v2.4.2` | `05b42c624433fc40578a4040d5cf5e36ddca8cde` |
| `anchore/sbom-action` | `v0` | `e22c389904149dbc22b58101806040fa8d37a610` |
| `CycloneDX/gh-gomod-generate-sbom` | `v2.0.0` | `efc74245d6802c8cefd925620515442756c70d8f` |
| `sigstore/cosign-installer` | `v4.1.0` | `ba7bc0a3fef59531c69a25acd34668d6d3fe6f22` |

## Maintenance Rule

When updating an action:

1. Resolve the tag to a commit SHA.
2. Use the SHA in workflow `uses:`.
3. Keep the tag in a comment.
4. Update this table.
5. Run `go test ./...` so `internal/assurance` checks the policy.
