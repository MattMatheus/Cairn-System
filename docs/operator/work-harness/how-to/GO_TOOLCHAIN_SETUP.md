# Go Toolchain Setup

## Purpose
Ensure environments can run `go version` and module-scoped `go test ./...` for Cairn Go components.

## Required Version
- Source of truth: `products/memory-cli/go.mod` and `products/tool-cli/go.mod`
- Current minimum: `go 1.22`

## Preflight
- `products/work-harness/tools/check_go_toolchain.sh`
- `cd products/tool-cli && go test ./...`
- `cd products/memory-cli && go test ./...`
