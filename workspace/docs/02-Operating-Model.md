# Operating Model

This page explains the model. For actual operator flow, use `workspace/docs/HUMANS.md` and `workspace/docs/08-Human-Runbook.md`.

## Planes
- Delivery Plane: planning -> architect -> engineering -> qa -> pm -> cycle-close
- Knowledge Plane: idea -> candidate -> active-test -> supported|falsified|parked

## High-Level Rules
1. No silent state changes.
2. No stage skipping.
3. Human approval is required for canonical promotion from staging/suggested outputs.
4. Claim promotion from active-test requires linked artifact evidence.
5. Safety decisions are mandatory for sensitive ingestion.

## Outcomes
- Delivery outputs are quality-gated and reviewable.
- Knowledge outputs are testable, traceable, and evidence-linked.
- Human and agent roles remain explicit.
