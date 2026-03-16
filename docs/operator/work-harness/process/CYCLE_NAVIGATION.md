---
title: Cycle Navigation
description: Minimal route through an work harness delivery cycle, from branch check to stage launch to observer closure.
doc_type: procedure
status: draft
domain: process
related:
  - [[STAGE_CONTRACTS]]
  - [[BRANCH_AND_COMMIT_INVARIANTS]]
---

# Cycle Navigation

## Summary

An work harness cycle starts with branch safety, moves through a stage launcher and backlog state transition, and ends only after observer output plus a single cycle commit.

## Intended Audience

- operators entering a cycle
- agents that need the shortest correct execution path

## Preconditions

- work is happening on the required branch, normally `dev`
- a stage needs to be launched or resumed

## Main Flow

1. Confirm the active branch matches `CAIRN_REQUIRED_BRANCH` (default `dev`).
2. Launch the needed stage with `products/work-harness/tools/launch_stage.sh <stage>`.
3. Follow the returned seed prompt and execute the current story or refinement step.
4. Move backlog state only through canonical lanes.
5. Before handoff or closure, run the verification expected by [[VERIFICATION_PATTERNS]].
6. Close the cycle with `products/work-harness/tools/run_observer_cycle.sh --cycle-id <cycle-id>`.
7. Create exactly one commit: `cycle-<cycle-id>`.

If engineering has no active stories, the correct outcome is `no stories` and the next action is PM refinement.

## Failure Modes

- stage launch on the wrong branch
- direct movement from `active` to `done`
- cycle completion without observer output
- intermediate commits before cycle closure

## References

- [[STAGE_CONTRACTS]]
- [[BRANCH_AND_COMMIT_INVARIANTS]]
- `docs/operator/work-harness/process/CYCLE_INDEX.md`
