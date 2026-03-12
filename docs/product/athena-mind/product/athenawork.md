# AthenaWork Product Guide

## What AthenaWork Is

AthenaWork is the operator-facing workflow system that constrains how agents execute work around AthenaMind.

It is designed so non-technical users can steer agents through explicit stages, queue policies, and evidence checkpoints.

## Why It Matters

- Provides stage guardrails (`planning`, `architect`, `engineering`, `qa`, `pm`, `cycle`).
- Enforces queue discipline and handoff structure.
- Captures observer evidence and cycle continuity.
- Reduces agent drift by routing execution through canonical prompts and specialist roles.

## Canonical Operator Assets

- Product operator guide: `products/athena-work/HUMANS.md`
- Product agent rules: `products/athena-work/AGENTS.md`
- Stage prompts: `products/athena-work/stage-prompts/active/`
- Specialist directory: `products/athena-work/staff-personas/STAFF_DIRECTORY.md`
- Queue system: `products/athena-work/delivery-backlog/`
- Work OS artifacts: `products/athena-work/operating-system/`
- Stage launchers and checks: `products/athena-work/tools/`
- Daily workspace surface: `workspace/`

## Typical Use Pattern

1. Start a stage with `./products/athena-work/tools/launch_stage.sh <stage>`.
2. Execute queue item under that stage prompt.
3. Validate with the relevant product checks.
4. Run observer: `./products/athena-work/tools/run_observer_cycle.sh --cycle-id <cycle-id>`.
5. Commit once per cycle.

## Related Docs

- `docs/product/athena-mind/product/athenawork-operator-reference.md`
- `docs/operator/athena-mind/getting-started/athenawork-quickstart.md`
- `docs/product/athena-mind/product/athenamind.md`
