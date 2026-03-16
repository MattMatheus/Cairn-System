# work harness Product Guide

## What work harness Is

work harness is the operator-facing workflow system that constrains how agents execute work around memory-cli.

It is designed so non-technical users can steer agents through explicit stages, queue policies, and evidence checkpoints.

## Why It Matters

- Provides stage guardrails (`planning`, `architect`, `engineering`, `qa`, `pm`, `cycle`).
- Enforces queue discipline and handoff structure.
- Captures observer evidence and cycle continuity.
- Reduces agent drift by routing execution through canonical prompts and specialist roles.

## Canonical Operator Assets

- Product operator guide: `products/work-harness/HUMANS.md`
- Product agent rules: `products/work-harness/AGENTS.md`
- Stage prompts: `products/work-harness/stage-prompts/active/`
- Specialist directory: `products/work-harness/staff-personas/STAFF_DIRECTORY.md`
- Queue system: `products/work-harness/delivery-backlog/`
- Work OS artifacts: `products/work-harness/operating-system/`
- Stage launchers and checks: `products/work-harness/tools/`
- Daily workspace surface: `workspace/`

## Typical Use Pattern

1. Start a stage with `./products/work-harness/tools/launch_stage.sh <stage>`.
2. Execute queue item under that stage prompt.
3. Validate with the relevant product checks.
4. Run observer: `./products/work-harness/tools/run_observer_cycle.sh --cycle-id <cycle-id>`.
5. Commit once per cycle.

## Related Docs

- `docs/product/memory-cli/product/work-harness-operator-reference.md`
- `docs/operator/memory-cli/getting-started/work-harness-quickstart.md`
- `docs/product/memory-cli/product/memorycli.md`
