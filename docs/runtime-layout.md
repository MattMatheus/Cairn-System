# Repo-Local Runtime Layout

## Purpose

This document defines the recommended local runtime boundary for Cairn repositories.

The default model is a repo-local `.athena/` directory that is ignored by git.

## Why

The shared mega-workspace pattern works for power users, but it creates avoidable confusion for most engineers:

- hidden cross-repo dependencies
- unclear ownership of local state
- path assumptions that do not transfer cleanly
- context pollution from unrelated notes, vaults, tools, and artifacts

Repo-local runtime state is easier to understand, easier to clean, and easier to automate.

## Default Contract

```text
<repo>/
  .athena/
    workspace/
    memory/
    artifacts/
    cache/
    runs/
    config/
```

## Directory Roles

### `.athena/workspace/`

Local AthenaWork runtime state:

- queue material
- generated handoffs
- transient working notes
- local copies of workspace artifacts when needed

### `.athena/memory/`

Local AthenaMind runtime state:

- sqlite or Mongo-backed memory roots
- embeddings
- indexes
- telemetry files
- snapshots

### `.athena/artifacts/`

Fetched bootstrap material:

- Azure blob artifacts
- Azure DevOps-delivered reference bundles
- tool payloads required for local setup

This keeps fetched material out of the committed tree.

### `.athena/cache/`

Disposable caches:

- temporary fetch products
- indexing caches
- short-lived local optimization outputs

### `.athena/runs/`

Generated outputs:

- observer outputs
- execution reports
- logs intended for local inspection

### `.athena/config/`

Machine-local overrides:

- local endpoint overrides
- backend selection
- operator-specific config not intended for git

## Committed vs Uncommitted Boundary

Keep committed:

- `SKILL.md`
- `AGENTS.md`
- product code
- platform docs
- templates
- bootstrap scripts
- small safe examples

Keep uncommitted:

- live task state
- memory stores
- fetched artifacts
- caches
- personal notes
- machine-local configuration

## Override Model

Default behavior should prefer the repo-local `.athena/` directory.

Advanced users may still opt into overrides:

1. repo-local `.athena/`
2. optional user-global Athena home
3. explicit environment variable overrides

Suggested environment shape:

- `ATHENA_HOME`
- `ATHENA_WORKSPACE_ROOT`
- `ATHENA_MEMORY_ROOT`

## Platform Direction

This runtime layout should be the default internal-beta posture.

The future bootstrap script should hydrate `.athena/artifacts/` and any required local runtime folders rather than assuming a shared global workspace.
