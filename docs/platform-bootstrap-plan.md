# AthenaPlatform Bootstrap Plan

## Objective

Create one practical platform from three current assets:

- `AthenaMind/`: memory engine
- `AthenaWork/`: staged delivery engine
- `work/AthenaWork/`: markdown work vault

## Proposed Platform Shape

### Layer 1: Runtime

Source today:

- `AthenaMind/cmd/memory-cli`
- `AthenaMind/internal/`
- `AthenaMind/memory/`

Responsibilities:

- memory write and retrieval
- indexing and evaluation
- snapshot and recovery
- telemetry and read gateway services

V1 runtime policy:

- default backend: `sqlite`
- content source: AthenaWork markdown
- integration surface: slim Go CLI
- optional advanced backend: document DB support for developers who want more power

### Layer 2: Work Control Plane

Source today:

- `AthenaWork/products/athena-work/`
- `AthenaWork/delivery-backlog/`
- `AthenaWork/operating-system/`
- `AthenaWork/stage-prompts/`
- `AthenaWork/tools/`

Responsibilities:

- stage launch and enforcement
- backlog state transitions
- observer cycle execution
- release and handoff contracts
- regression checks for workflow integrity

### Layer 3: Workspace Plane

Source today:

- `work/AthenaWork/00-Docs/`
- `work/AthenaWork/10 Work/`
- `work/AthenaWork/30 AI/`
- `work/AthenaWork/40 Assets/`
- `work/AthenaWork/70 Agents/`
- `work/AthenaWork/80 Research/`

Responsibilities:

- daily work capture
- project and area notes
- research claims and artifacts
- agent queue handling
- human-readable maps and dashboards

## Initial Canonical Direction

The most defensible short-term architecture is:

- keep AthenaMind as the runtime subsystem
- keep AthenaWork root as the executable workflow subsystem
- use the vault-style AthenaWork as the target workspace UX and content model

That gives the unified platform:

- executable automation
- durable governance
- a work surface that people can actually use daily

## Root-Level Scaffold Decision

For now, the root repo should own:

- platform-level docs
- product-owned destination directories
- shared tool directories
- workspace destination directories
- comparison and migration planning
- future shared contracts between runtime and workspace

It should not yet rewrite the embedded repositories in-place beyond creating the destination scaffold.

## Suggested Next Milestones

1. Define a canonical platform directory model for imported code versus live user content.
2. Normalize metadata contracts between delivery items and workspace notes.
3. Design how AthenaMind retrieval attaches to AthenaWork queues, runs, claims, and artifacts.
4. Decide whether user content lives inside the platform repo or in a separate private vault connected through contracts.
5. Import or mirror selected components into root-owned modules once the target layout is stable.

## Candidate Future Root Layout

```text
AthenaPlatform/
  README.md
  docs/
  products/
    athena-mind/
    athena-work/
  tools/
    platform/
    dev/
  workspace/
    docs/
    work/
    agents/
    research/
```

## Risks To Manage

- Embedded git repositories complicate a future single-repo import.
- The work vault contains live content and may require privacy boundaries.
- State models overlap but are not yet identical.
- Tooling expects different canonical paths across the current versions.
- Optional advanced backends should not complicate the default developer setup.

## Recommended Next Task

Produce a canonical target tree for `AthenaPlatform` and explicitly map each current directory into:

- keep as-is
- import
- merge
- archive
- ignore

