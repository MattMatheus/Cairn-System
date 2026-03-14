---
title: Go
description: Official Go documentation record for language, module, toolchain, and test-command semantics used by AthenaWork.
doc_type: canonical_source
status: draft
tool: go
vendor: Go Project
version_policy: repo_pinned
version: ""
docs_base: https://go.dev/doc/
allowed_domains:
  - go.dev
entrypoints:
  - https://go.dev/doc/
  - https://pkg.go.dev/cmd/go
consult_when:
  - checking go command behavior
  - reviewing module, workspace, or toolchain semantics
  - validating test, build, or package behavior in Go code
avoid_when:
  - the task is only about AthenaWork process
notes:
  - Prefer the repo's declared Go version when available, then use the matching official Go docs.
---

# Go

## Summary

Use the official Go docs for language and `go` tool behavior. When a repository pins a Go version, interpret command and toolchain behavior in that version's context.

## References

- https://go.dev/doc/
- https://pkg.go.dev/cmd/go
