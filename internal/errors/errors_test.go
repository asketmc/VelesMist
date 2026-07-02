// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package errors

import "testing"

func TestExitCodeMapping(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "nil", err: nil, want: ExitSuccess},
		{name: "invalid", err: New(InvalidInput, "bad input"), want: ExitInvalidInput},
		{name: "upstream", err: New(Upstream, "steam down"), want: ExitUpstream},
		{name: "rate limited", err: New(RateLimited, "slow down"), want: ExitUpstream},
		{name: "timeout", err: New(NetworkTimeout, "timeout"), want: ExitUpstream},
		{name: "internal", err: New(Internal, "bug"), want: ExitInternal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCode(tt.err); got != tt.want {
				t.Fatalf("ExitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}
