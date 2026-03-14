# AthenaPlatform

External brand direction: `Cairn`

AthenaPlatform is the current internal platform workspace for the broader `Cairn` direction.

Today it contains the working Athena product family:

- `AthenaMind`: local developer memory and retrieval
- `AthenaWork`: staged workflow and delivery operations
- `AthenaUse`: governed tool discovery and context emission
- `workspace/`: markdown-native operating surface

## Current State

This root repository was initialized on 2026-03-12 to unify AthenaMind and AthenaWork into a single usable platform.

## North Star

`Cairn` is being built as:

- a disciplined personal PKM and human-agent work system now
- a future offline-friendly, curated knowledge tool for remote ministry contexts later

That means the platform should prefer local-first operation, explicit review, low dependency weight, and outputs that remain useful under weak connectivity.

## Naming Posture

Current posture:

- `Cairn` is the external umbrella brand and long-term outward-facing name
- `AthenaPlatform`, `AthenaMind`, `AthenaUse`, and `AthenaWork` remain the current internal working names in this repository

Rename policy:

- avoid disruptive broad renames while product boundaries are still moving
- prefer a later coordinated rename pass once command surfaces, repos, and docs are stable enough to change together

Today the platform contains:

- `docs/` for root-level platform documentation
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
- `sqlite` as the storage path
- a Go CLI as the main integration surface

### `products/athena-work`

A unified work operating system combining:

- staged delivery and validation tooling
- markdown-native workspace structure
- reusable instructions, prompts, templates, and contracts
