# Operating System

Supporting process-improvement and control-plane material for AthenaWork.

This is not the recommended first-stop operator entrypoint for internal beta users.

## Use First
- `products/athena-work/HUMANS.md`
- `products/athena-work/DEVELOPMENT_CYCLE.md`
- `workspace/docs/HUMANS.md`
- `workspace/docs/08-Human-Runbook.md`

## Purpose
Treat process quality as a product: measurable, iterative, and testable.

## Core Loop
1. Capture process gaps in `operating-system/backlog/intake/`.
2. Refine and rank into `operating-system/backlog/active/`.
3. Implement process changes.
4. QA the process change impact.
5. Run Observer at cycle boundary and capture metadata deltas.
6. Commit once per cycle with observer report + artifacts.
7. Measure KPI deltas in `operating-system/metrics/`.
8. Keep/adjust/revert via decision records.
9. Maintain explicit shipped checkpoint bundles before declaring release completion.

## Documentation Maintenance Rule
Any accepted process change must include documentation sync:
- `HUMANS.md` (operator quick guide)
- `AGENTS.md` (agent discovery/operating guide)
- `DEVELOPMENT_CYCLE.md` (canonical stage behavior)
- relevant `workspace/docs/*` guidance when the user-facing workflow changes

## Structure
- `delivery-backlog/`: state-machine for process-improvement stories
- `experiments/`: hypothesis-driven trials
- `decisions/`: process ADR-style records
- `metrics/`: process KPIs and trend snapshots
- `playbooks/`: operational runbooks for cycle stages
- `retros/`: sprint and incident retrospectives
- `handoff/`: current process handoff status
- `observer/`: cycle-boundary observer reports and templates

## Status
- `backlog/`, `observer/`, and selected `handoff/` artifacts can affect active workflow.
- `decisions/`, `experiments/`, `metrics/`, `playbooks/`, and `retros/` are mainly supporting/reference surfaces unless a current change explicitly targets them.
