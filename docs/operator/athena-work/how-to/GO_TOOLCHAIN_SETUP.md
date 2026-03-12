# Go Toolchain Setup

## Purpose
Ensure environments can run `go version` and `go test ./...` for state-harness Go components.

## Required Version
- Source of truth: `go.mod`
- Current minimum: `go 1.22`

## Preflight
- `products/athena-work/tools/check_go_toolchain.sh`
- `go test ./...`
