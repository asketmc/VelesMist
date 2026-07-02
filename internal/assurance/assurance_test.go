// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package assurance

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var usesPattern = regexp.MustCompile(`(?m)^\s*uses:\s+(.+?)\s*$`)
var pinnedPattern = regexp.MustCompile(`(?m)^\s*uses:\s+[^@\s]+@[0-9a-f]{40}(?:\s+#\s+.+)?\s*$`)

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
