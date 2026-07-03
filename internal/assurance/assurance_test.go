// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package assurance

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

var usesPattern = regexp.MustCompile(`(?m)^\s*uses:\s+(.+?)\s*$`)
var pinnedPattern = regexp.MustCompile(`(?m)^\s*uses:\s+[^@\s]+@[0-9a-f]{40}(?:\s+#\s+.+)?\s*$`)
var pinnedUsePattern = regexp.MustCompile(`(?m)^\s*uses:\s+([^@\s]+)@([0-9a-f]{40})(?:\s+#\s+(\S+))?\s*$`)
var pinningDocRowPattern = regexp.MustCompile("(?m)^\\|\\s+`([^`]+)`\\s+\\|\\s+`([^`]+)`\\s+\\|\\s+`([0-9a-f]{40})`\\s+\\|$")

func TestWorkflowsUsePinnedActions(t *testing.T) {
	root := repoRoot(t)
	paths, err := filepath.Glob(filepath.Join(root, ".github", "workflows", "*.yml"))
	if err != nil {
		t.Fatalf("glob workflows: %v", err)
	}
	var violations []string
	for _, path := range paths {
		body := readFile(t, path)
		for _, match := range usesPattern.FindAllStringSubmatch(body, -1) {
			if strings.HasPrefix(match[1], "./") || strings.HasPrefix(match[1], "docker://") {
				continue
			}
			if !pinnedPattern.MatchString(match[0]) {
				violations = append(violations, filepath.Base(path)+": "+match[1])
			}
		}
	}
	if len(violations) > 0 {
		t.Fatalf("unpinned workflow actions: %v", violations)
	}
}

func TestWorkflowsDeclarePermissions(t *testing.T) {
	root := repoRoot(t)
	paths, err := filepath.Glob(filepath.Join(root, ".github", "workflows", "*.yml"))
	if err != nil {
		t.Fatalf("glob workflows: %v", err)
	}
	for _, path := range paths {
		body := readFile(t, path)
		if !strings.Contains(body, "\npermissions:\n") {
			t.Fatalf("%s has no top-level permissions block", filepath.Base(path))
		}
	}
}

func TestWorkflowPinningDocsMatchWorkflows(t *testing.T) {
	root := repoRoot(t)
	workflowPins := workflowPins(t, root)
	docPins := pinningDocPins(t, filepath.Join(root, "docs", "WORKFLOW_PINNING.md"))

	var missing []string
	for action, pin := range workflowPins {
		docPin, ok := docPins[action]
		if !ok {
			missing = append(missing, action)
			continue
		}
		if docPin.sha != pin.sha || docPin.tag != pin.tag {
			t.Fatalf("docs/WORKFLOW_PINNING.md has stale pin for %s: got %s/%s, want %s/%s", action, docPin.tag, docPin.sha, pin.tag, pin.sha)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf("docs/WORKFLOW_PINNING.md missing workflow actions: %v", missing)
	}
	var extra []string
	for action := range docPins {
		if _, ok := workflowPins[action]; !ok {
			extra = append(extra, action)
		}
	}
	if len(extra) > 0 {
		sort.Strings(extra)
		t.Fatalf("docs/WORKFLOW_PINNING.md lists unused workflow actions: %v", extra)
	}
}

func TestCodeownersCoversSensitivePaths(t *testing.T) {
	root := repoRoot(t)
	body := readFile(t, filepath.Join(root, ".github", "CODEOWNERS"))
	required := []string{
		".github/workflows/*",
		"go.mod",
		"go.sum",
		"scripts/*",
		"SECURITY.md",
		"docs/VERIFY_RELEASE.md",
		"docs/RELEASE_SECURITY.md",
		"docs/THREAT_MODEL.md",
	}
	for _, pattern := range required {
		if !strings.Contains(body, pattern) {
			t.Fatalf("CODEOWNERS missing %s", pattern)
		}
	}
}

func TestPullRequestTemplateCoversQAAndSecurity(t *testing.T) {
	root := repoRoot(t)
	body := readFile(t, filepath.Join(root, ".github", "pull_request_template.md"))
	requiredSections := []string{
		"## Summary",
		"## Behavior changed?",
		"## Output contract changed?",
		"## Security/privacy impact?",
		"## Tests added/updated?",
		"## Local Validation Evidence",
		"## Artifact/docs impact",
		"## Breaking changes?",
		"## Non-Goals",
	}
	for _, section := range requiredSections {
		if !strings.Contains(body, section) {
			t.Fatalf("pull request template missing section %q", section)
		}
	}

	requiredChecklist := []string{
		"No secrets/tokens/cookies/API keys committed",
		"No real Steam credentials in tests",
		"No real network dependency in unit tests",
		"Output contract updated if JSON/table changed",
		"README/docs updated if CLI changed",
		"`make verify` run locally",
	}
	for _, item := range requiredChecklist {
		if !strings.Contains(body, item) {
			t.Fatalf("pull request template missing checklist item %q", item)
		}
	}
}

func TestOSSAssuranceCoversRequiredControls(t *testing.T) {
	root := repoRoot(t)
	body := readFile(t, filepath.Join(root, "docs", "OSS_ASSURANCE.md"))
	lowerBody := strings.ToLower(body)
	forbidden := []string{"certified", "certification", "external audit"}
	for _, term := range forbidden {
		if strings.Contains(lowerBody, term) {
			t.Fatalf("docs/OSS_ASSURANCE.md contains external assurance wording %q", term)
		}
	}

	allowedStatuses := map[string]bool{
		"implemented":             true,
		"partial":                 true,
		"planned":                 true,
		"not applicable":          true,
		"requires GitHub setting": true,
		"requires first release":  true,
	}
	controls := ossAssuranceControls(t, body)
	for control, status := range controls {
		if !allowedStatuses[status] {
			t.Fatalf("docs/OSS_ASSURANCE.md control %q has unsupported status %q", control, status)
		}
	}

	requiredControls := []string{
		"OpenSSF Scorecard",
		"CI",
		"CodeQL",
		"Dependabot",
		"Dependency Review",
		"Secret scanning / push protection",
		"REUSE",
		"SPDX SBOM",
		"CycloneDX SBOM",
		"SLSA / GitHub artifact attestations",
		"Sigstore / cosign",
		"govulncheck",
		"OSV Scanner",
		"Semgrep",
		"workflow pinning",
		"Security Insights",
		"CODEOWNERS",
		"branch protection",
		"QA_MAP",
		"ARTIFACTS",
		"output schemas/contracts",
		"release verification docs",
	}
	for _, control := range requiredControls {
		if _, ok := controls[control]; !ok {
			t.Fatalf("docs/OSS_ASSURANCE.md missing required control %q", control)
		}
	}
}

type actionPin struct {
	tag string
	sha string
}

func workflowPins(t *testing.T, root string) map[string]actionPin {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(root, ".github", "workflows", "*.yml"))
	if err != nil {
		t.Fatalf("glob workflows: %v", err)
	}
	pins := map[string]actionPin{}
	for _, path := range paths {
		body := readFile(t, path)
		for _, match := range pinnedUsePattern.FindAllStringSubmatch(body, -1) {
			action := match[1]
			if strings.HasPrefix(action, "./") || strings.HasPrefix(action, "docker://") {
				continue
			}
			pins[action] = actionPin{tag: match[3], sha: match[2]}
		}
	}
	return pins
}

func pinningDocPins(t *testing.T, path string) map[string]actionPin {
	t.Helper()
	body := readFile(t, path)
	pins := map[string]actionPin{}
	for _, match := range pinningDocRowPattern.FindAllStringSubmatch(body, -1) {
		pins[match[1]] = actionPin{tag: match[2], sha: match[3]}
	}
	return pins
}

func ossAssuranceControls(t *testing.T, body string) map[string]string {
	t.Helper()
	controls := map[string]string{}
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") || strings.Contains(line, "---") {
			continue
		}
		columns := strings.Split(strings.Trim(line, "|"), "|")
		if len(columns) != 3 {
			continue
		}
		control := strings.TrimSpace(columns[0])
		status := strings.TrimSpace(columns[1])
		if control == "Control" && status == "Status" {
			continue
		}
		controls[control] = status
	}
	if len(controls) == 0 {
		t.Fatal("docs/OSS_ASSURANCE.md controls table was not parsed")
	}
	return controls
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(body)
}
