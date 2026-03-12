# AthenaWork Operator Reference

## Purpose

Operational reference for running AthenaWork directly in this repository.

## Operator Root Paths

- `products/athena-work/HUMANS.md`
- `products/athena-work/DEVELOPMENT_CYCLE.md`
- `products/athena-work/tools/`
- `products/athena-work/stage-prompts/active/`
- `products/athena-work/delivery-backlog/`
- `products/athena-work/operating-system/`
- `products/athena-work/staff-personas/`

## Stage Launch Commands

```bash
./products/athena-work/tools/launch_stage.sh planning
./products/athena-work/tools/launch_stage.sh architect
./products/athena-work/tools/launch_stage.sh engineering
./products/athena-work/tools/launch_stage.sh qa
./products/athena-work/tools/launch_stage.sh pm
./products/athena-work/tools/launch_stage.sh cycle
```

Observer:

```bash
./products/athena-work/tools/run_observer_cycle.sh --cycle-id <cycle-id>
```

## Required Execution Rules

- Respect lane boundaries (`engineering` vs `architecture`).
- Run product-specific validation before handoff.
- Run observer after completed cycles.
- Commit exactly once per cycle (`cycle-<cycle-id>`).
- Treat `done` as QA-complete, not auto-shipped.

## Human/Non-Technical Steering Model

- Non-technical operators choose stage + goal.
- Stage prompt and specialist persona constrain agent behavior.
- Queue artifacts provide explicit acceptance criteria.
- Observer and release-handoff artifacts provide auditability.

## Related

- `docs/product/athena-mind/product/athenawork.md`
- `products/athena-work/DEVELOPMENT_CYCLE.md`
- `docs/product/athena-mind/product/athenamind.md`
