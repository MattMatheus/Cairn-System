# System Overview

For day-to-day usage, start with `workspace/docs/HUMANS.md` and `workspace/docs/08-Human-Runbook.md`.

AthenaWork uses a dual-plane management model:
- Delivery Plane: deterministic work execution from intake to done.
- Knowledge Plane: claim validation lifecycle from idea to supported or falsified.

## Design Goals
- Human-readable and agent-operable.
- Deterministic state movement.
- Evidence-backed knowledge promotion.
- Safety-first ingestion and promotion controls.

## Control Surfaces
- Human entrypoint: `workspace/docs/HUMANS.md`
- Agent entrypoint: `workspace/AGENTS.md`
- Operating contracts: `products/athena-work/operating-system-vault/*`
- Delivery queue lanes: `workspace/agents/delivery/*`
- Claim graph zones: `workspace/research/*`
