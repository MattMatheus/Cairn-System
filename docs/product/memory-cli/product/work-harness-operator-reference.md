# work harness Operator Reference

## Purpose

Operational reference for running work harness directly in this repository.

## Operator Root Paths

- `products/work-harness/HUMANS.md`
- `products/work-harness/DEVELOPMENT_CYCLE.md`
- `products/work-harness/tools/`
- `products/work-harness/stage-prompts/active/`
- `products/work-harness/delivery-backlog/`
- `products/work-harness/operating-system/`
- `products/work-harness/staff-personas/`

## Stage Launch Commands

```bash
./products/work-harness/tools/launch_stage.sh planning
./products/work-harness/tools/launch_stage.sh architect
./products/work-harness/tools/launch_stage.sh engineering
./products/work-harness/tools/launch_stage.sh qa
./products/work-harness/tools/launch_stage.sh pm
./products/work-harness/tools/launch_stage.sh cycle
```

Observer:

```bash
./products/work-harness/tools/run_observer_cycle.sh --cycle-id <cycle-id>
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

- `docs/product/memory-cli/product/work-harness.md`
- `products/work-harness/DEVELOPMENT_CYCLE.md`
- `docs/product/memory-cli/product/memorycli.md`
