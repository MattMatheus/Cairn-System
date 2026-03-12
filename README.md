# AthenaPlatform

AthenaPlatform is the unified internal platform for:

- `AthenaMind`: local developer memory and retrieval
- `AthenaWork`: staged workflow and delivery operations
- `AthenaUse`: governed tool discovery and context emission
- `workspace/`: markdown-native operating surface

## Current State

This root repository was initialized on 2026-03-12 to unify AthenaMind and AthenaWork into a single usable platform.

Today the platform contains:

- `docs/` for root-level platform analysis and integration planning
- `products/` as the canonical home for platform-owned product code
- `tools/` as the shared tooling surface
- `workspace/` as the canonical markdown workspace surface
- `.athena/` as the intended repo-local runtime area for uncommitted operational state

## Start Here

- Internal beta start: `INTERNAL_BETA.md`
- Unified quickstart: `PLATFORM_QUICKSTART.md`
- Local bootstrap: `tools/dev/bootstrap_platform.sh`
- Operator docs: `docs/operator/README.md`
- Product docs: `docs/product/README.md`
- Platform scaffold: `docs/platform-layout.md`
- Repo-local runtime contract: `docs/runtime-layout.md`

Historical planning and migration material remains available in:

- `docs/migration/`
- `docs/athenawork-versions-diff.md`
- `docs/platform-bootstrap-plan.md`

## Platform Layout

```text
AthenaPlatform/
  .athena/
  docs/
  products/
    athena-mind/
    athena-work/
    athena-use/
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

Uncommitted runtime state should live under a repo-local `.athena/` directory so each repository has an isolated Athena operating area without leaking personal state or generated artifacts upstream.

That runtime area is intended for:

- local AthenaWork workspace state
- AthenaMind memory roots
- fetched bootstrap artifacts
- caches
- run outputs
- machine-local config

## Product Intent

### `products/athena-mind`

A slim developer memory system with:

- markdown-first ingestion
- `sqlite` as the default local backend
- optional document database support for advanced users
- a Go CLI as the main integration surface

### `products/athena-work`

A unified work operating system combining:

- staged delivery and validation tooling
- markdown-native workspace structure
- reusable instructions, prompts, templates, and contracts
