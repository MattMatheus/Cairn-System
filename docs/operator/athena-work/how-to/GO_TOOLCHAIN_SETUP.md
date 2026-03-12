# Go Toolchain Setup

## Purpose
Ensure environments can run `go version` and module-scoped `go test ./...` for AthenaPlatform Go components.

## Required Version
- Source of truth: `products/athena-mind/go.mod` and `products/athena-use/go.mod`
- Current minimum: `go 1.22`

## Preflight
- `products/athena-work/tools/check_go_toolchain.sh`
- `cd products/athena-use && go test ./...`
- `cd products/athena-mind && go test ./...`
