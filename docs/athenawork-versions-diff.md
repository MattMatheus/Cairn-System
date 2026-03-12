# AthenaWork Version Comparison

## Compared Directories

- `AthenaWork/`
- `work/AthenaWork/`

## Executive Summary

These are not minor variants. They represent two different product shapes:

- `AthenaWork/` is a productized repository centered on staged delivery orchestration, backlog flow, launch tooling, CI, and testable operating-system contracts.
- `work/AthenaWork/` is a markdown vault centered on real work execution, research capture, agent queues, evidence-backed claim handling, and write-boundary governance.

The practical unification target is not choosing one over the other. It is combining:

- the runtime discipline and tooling from `AthenaWork/`
- the usable work-vault information architecture from `work/AthenaWork/`

## Structural Differences

### `AthenaWork/` adds

- product packaging under `products/athena-work/`
- delivery backlog lanes under `delivery-backlog/`
- operating-system docs under `operating-system/`
- `stage-prompts/`
- `staff-personas/`
- CI and release artifacts such as `.github/`, `azure-pipelines.yml`, and docs website assets
- extensive shell tooling and regression tests in `tools/`
- `knowledge-base/` and `product-research/`

### `work/AthenaWork/` adds

- a vault-style knowledge and work layout:
  - `00-Docs/`
  - `10 Work/`
  - `30 AI/`
  - `40 Assets/`
  - `70 Agents/`
  - `80 Research/`
- live work and research content
- explicit write-boundary model for agents vs humans
- stricter frontmatter and metadata contracts for operational notes
- escape-path handling for blocked or ambiguous work
- vault-specific utilities such as `tools/prune_user_content.sh`

## Behavioral Differences

### Root `AthenaWork/`

- Agent flow is stage-launch oriented.
- The first-run path is based on `DEVELOPMENT_CYCLE.md`, active backlog lanes, and shell launch scripts.
- Branch policy is explicit and repo-centric, with `dev` as the enforced working branch.
- Commit policy is cycle-based with observer runs after each completed cycle.

### `work/AthenaWork/`

- Agent flow is note-centric and metadata-driven.
- The first-run path is based on docs, operating-system contracts, maps, and templates.
- Governance is framed around note state transitions, evidence linking, and zone-based write permissions.
- Safety and escalation behavior are more explicit for ambiguous or sensitive work.

## Document-Level Drift

## `AGENTS.md`

- `AthenaWork/AGENTS.md` is oriented around staged delivery and product operations.
- `work/AthenaWork/AGENTS.md` is oriented around a supervised claim-first dual-plane workspace.

## `README.md`

- `AthenaWork/README.md` is a slim product-distribution entrypoint.
- `work/AthenaWork/README.md` is a full operating-system overview for a live workspace.

## `HUMANS.md`

- `AthenaWork/HUMANS.md` is a procedural operator guide focused on branch, stage, and observer workflow.
- `work/AthenaWork/HUMANS.md` behaves more like a vault home page linking docs, templates, maps, work, and research dashboards.

## Branch Context

At comparison time on 2026-03-12:

- `AthenaWork/` branch: `dev`
- `work/AthenaWork/` branch: `main`

That likely explains part of the divergence, but the directory and documentation differences are large enough that branch drift alone does not explain the whole split.

## Recommended Merge Direction

Use the following merge posture:

1. Treat `AthenaWork/` as the source for executable workflow machinery.
2. Treat `work/AthenaWork/` as the source for workspace information architecture and human-usable note structure.
3. Define a root platform contract that allows AthenaMind to supply memory, retrieval, evaluation, and telemetry to the unified work system.

## Immediate Integration Targets

- Unify state models across backlog lanes and vault queues.
- Align metadata/frontmatter contracts with delivery stages.
- Map observer outputs to workspace artifacts.
- Decide whether `products/athena-work/` becomes the canonical app surface or whether the vault layout becomes the canonical workspace surface.
- Separate reusable system contracts from live user content.

