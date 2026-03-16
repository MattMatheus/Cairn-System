# Program Control Plane

## Core Rules
1. Delivery and knowledge work run in parallel planes with explicit handoffs.
2. No silent state transitions.
3. Human approval is required for high-impact promotions.
4. Every accepted claim state requires evidence linkage.
5. Safety policy is a hard gate for ingestion and promotion.

## Delivery Plane Contract
- Stages: `planning -> architect -> engineering -> qa -> pm -> cycle-close`
- Commit discipline: one commit per completed cycle.
- Stage skip is disallowed.

## Knowledge Plane Contract
- States: `idea -> candidate -> active-test -> supported|falsified|parked`
- `candidate -> active-test` requires test plan.
- `active-test -> supported|falsified` requires one or more linked artifacts.

## Human + Agent Boundary
- Agent Write zone: queue, runs, memory candidates, staging.
- Agent Suggest zone: canonical work/research maps and notes.
- Human decision required for canonical promotion from staging/suggested outputs.

## Traceability Requirements
- Delivery item -> references claim(s)/artifact(s) when behavior changes.
- Claim -> references artifact evidence and related delivery work when applicable.
