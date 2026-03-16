# Cairn

This repository is the unified Cairn platform workspace.

Today it contains the working product surfaces:

- `memory-cli`: local developer memory and retrieval
- `work harness`: staged workflow and delivery operations
- `tool-cli`: governed tool discovery and context emission
- `workspace/`: markdown-native operating surface

## Current State

This root repository was initialized on 2026-03-12 to unify memory-cli and work harness into a single usable platform.

## North Star

`Cairn` is being built as:

- a disciplined personal PKM and human-agent work system now
- a future offline-friendly, curated knowledge tool for remote ministry contexts later

That means the platform should prefer local-first operation, explicit review, low dependency weight, and outputs that remain useful under weak connectivity.

## Naming Posture

Current posture:

- the platform umbrella is `Cairn`
- `memory-cli`, `tool-cli`, and `work harness` remain the current internal product names in this repository
- the GitHub repository name `cairn-system` is reserved to isolate the UI component, not to rename the platform umbrella again

Rename policy:

- keep the platform umbrella on `Cairn`
- defer product-family renames until command surfaces, binaries, and docs are stable enough to change together

Today the platform contains:

- `docs/` for root-level platform documentation
- `products/` as the canonical home for platform-owned product code
- `tools/` as the shared tooling surface
- `workspace/` as the canonical markdown workspace surface
- `.cairn/` as the intended repo-local runtime area for uncommitted operational state

## Start Here

- Internal beta start: `INTERNAL_BETA.md`
- Unified quickstart: `PLATFORM_QUICKSTART.md`
- Local bootstrap: `tools/dev/bootstrap_platform.sh`
- Operator docs: `docs/operator/README.md`
- Product docs: `docs/product/README.md`
- Platform scaffold: `docs/platform-layout.md`
- Repo-local runtime contract: `docs/runtime-layout.md`

## Platform Layout

```text
Cairn/
  .cairn/
  docs/
  products/
    memory-cli/
    work-harness/
    tool-cli/
  tools/
    platform/
    dev/
  workspace/
    docs/
    work/
    agents/
    research/
```

## Runtime Model

Committed product and platform inputs stay in the repository.

Uncommitted runtime state should live under a repo-local `.cairn/` directory so each repository has an isolated Cairn operating area without leaking personal state or generated artifacts upstream.

That runtime area is intended for:

- local work harness workspace state
- memory-cli memory roots
- fetched bootstrap artifacts
- caches
- run outputs
- machine-local config

## Product Intent

### `products/memory-cli`

A slim developer memory system with:

- markdown-first ingestion
- `sqlite` as the storage path
- a Go CLI as the main integration surface

### `products/work-harness`

A unified work operating system combining:

- staged delivery and validation tooling
- markdown-native workspace structure
- reusable instructions, prompts, templates, and contracts
