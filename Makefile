# SPDX-FileCopyrightText: 2026 VelesMist contributors
# SPDX-License-Identifier: MIT

APP := velesmist
PKG := github.com/asketmc/VelesMist
DIST := dist
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo unknown)
DIRTY ?= $(shell test -n "$$(git status --porcelain 2>/dev/null)" && echo dirty || echo clean)
LDFLAGS := -s -w -X $(PKG)/internal/version.Version=$(VERSION) -X $(PKG)/internal/version.Commit=$(COMMIT) -X $(PKG)/internal/version.BuildDate=$(BUILD_DATE) -X $(PKG)/internal/version.Dirty=$(DIRTY)
GOVULNCHECK_VERSION ?= latest

.PHONY: test lint vet vuln build snapshot-release clean

test:
	go test ./...

lint:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

vet:
	go vet ./...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(DIST)/$(APP) ./cmd/velesmist

snapshot-release: clean
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(DIST)/$(APP)-linux-amd64 ./cmd/velesmist
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(DIST)/$(APP)-linux-arm64 ./cmd/velesmist
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(DIST)/$(APP)-windows-amd64.exe ./cmd/velesmist
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(DIST)/$(APP)-windows-arm64.exe ./cmd/velesmist
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(DIST)/$(APP)-darwin-amd64 ./cmd/velesmist
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(DIST)/$(APP)-darwin-arm64 ./cmd/velesmist

clean:
	rm -rf $(DIST)
