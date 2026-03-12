# AthenaPlatform Layout

## Purpose

This document defines the root scaffold for the unified platform.

The platform must support:

- multiple products
- shared tooling
- a canonical markdown workspace
- documentation that is separate from product code

## Canonical Root Tree

```text
AthenaPlatform/
  .athena/
    workspace/
    memory/
    artifacts/
    cache/
    runs/
    config/
  docs/
    decisions/
    migration/
    operator/
    product/
  products/
    athena-mind/
    athena-work/
    athena-use/
  tools/
    platform/
    dev/
    migration/
  workspace/
    docs/
    work/
    agents/
    research/
    templates/
```

## Responsibilities

### `docs/`

Platform-owned documentation:

- architecture and decisions
- migration plans
- operator guidance
- product overviews

### `products/athena-mind/`

Canonical home for the slim memory product:

- markdown ingestion
- local retrieval
- `sqlite` default backend
- optional document DB backend for advanced users
- CLI-first integration

### `products/athena-work/`

Canonical home for the unified work system:

- workflow contracts
- launch scripts
- prompts
- validations
- stage and observer machinery

### `products/athena-use/`

Canonical home for the governed tool-context product:

- approved tool registry
- tool discovery and context emission
- tool-validation logic
- OpenTelemetry-instrumented tool selection and usage tracing

### `tools/`

Shared platform tooling, separated from product-specific tools:

- `platform/`: root-level orchestration helpers
- `dev/`: local developer helpers
- `migration/`: import, diff, cleanup, and classification helpers

### `workspace/`

Canonical markdown work surface:

- `docs/`: workspace-native docs
- `work/`: project and area notes
- `agents/`: queue, runs, staging, and control surfaces
- `research/`: claims, concepts, artifacts, and maps
- `templates/`: note and task templates

### `.athena/`

Repo-local runtime area, intentionally excluded from version control:

- `workspace/`: local AthenaWork operational state
- `memory/`: AthenaMind memory roots and backend state
- `artifacts/`: fetched bootstrap assets from Azure or other approved sources
- `cache/`: disposable local caches
- `runs/`: generated run output and transient logs
- `config/`: machine-local overrides

## Design Rules

1. Product code lives under `products/`.
2. Shared scripts live under `tools/`.
3. Human-facing markdown workspace content lives under `workspace/`.
4. Platform documentation lives under `docs/` rather than being mixed into product folders.
5. Optional advanced backends must not increase the complexity of the default local path.
6. Repo-local operational state belongs in `.athena/`, not in committed product folders.
