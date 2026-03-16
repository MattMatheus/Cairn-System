# Delivery Workflow

This page describes the delivery flow. For the shortest active path, pair it with `products/work-harness/HUMANS.md`.

## Stage Sequence
1. Planning
2. Architect
3. Engineering
4. QA
5. PM
6. Cycle Closure

## Lane Model
- Architecture: `intake -> active -> qa -> done`
- Engineering: `intake -> active -> qa -> done`

## Required Exit Conditions
- Planning: session captured, intake created, next stage named.
- Architect: decision scope explicit, outputs mapped to downstream implementation.
- Engineering: acceptance criteria complete, tests updated and passing, handoff complete.
- QA: acceptance and regression evidence explicit, transition explicit.
- PM: intake validated, active queue ranked, control-plane synchronized.
- Cycle closure: observer/report artifact produced, single cycle commit discipline preserved.

## Defect Rule
- QA defects return work to `active` with linked bug/task.
