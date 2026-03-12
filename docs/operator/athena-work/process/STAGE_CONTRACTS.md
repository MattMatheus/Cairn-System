---
title: Stage Contracts
description: Short contract per stage describing the outputs and gates required before AthenaWork can hand off or transition state.
doc_type: claim
status: draft
domain: process
related:
  - [[HANDOFF_EXPECTATIONS]]
  - [[VERIFICATION_PATTERNS]]
---

# Stage Contracts

## Summary

Each AthenaWork stage has a contract: required outputs, required evidence, and an allowed state transition. A stage is not done because time was spent; it is done because its contract is satisfied.

## Intended Audience

- stage agents
- operators reviewing whether a handoff is legitimate

## Preconditions

- a story or planning artifact is currently in a named stage

## Main Flow

Use these minimum contracts:

- Planning: finalized planning artifact, explicit next-stage recommendation, intake artifacts created
- Architect: scoped architecture output package, validation commands recorded, follow-on implementation paths listed
- PM: intake validation passes, active queue ranked, program board updated
- Engineering: acceptance criteria implemented, touched behavior verified, handoff package complete, story moves `active -> qa`
- QA: evidence recorded, regression gate passes, QA result artifact exists, state moves `qa -> done` or `qa -> active`
- Cycle closure: observer report exists and exactly one cycle commit is made

When in doubt, check whether the stage can produce explicit evidence for the next stage to act without re-deriving context.

## Details

Stage contracts answer:

- what artifact must exist
- what verification must pass
- what state transition is legal
- what the next stage should be able to trust

## Failure Modes

- stage output without explicit next-stage usefulness
- missing artifact or missing evidence
- state transition that bypasses `qa`
- stage marked done while risks/questions remain implicit

## References

- [[HANDOFF_EXPECTATIONS]]
- [[VERIFICATION_PATTERNS]]
- `docs/operator/athena-work/process/STAGE_EXIT_GATES.md`
