---
title: Branch And Commit Invariants
description: Non-negotiable work harness rules for branch discipline, state movement, observer closure, and single-commit cycles.
doc_type: claim
status: draft
domain: process
related:
  - [[CYCLE_NAVIGATION]]
  - [[STAGE_CONTRACTS]]
---

# Branch And Commit Invariants

## Summary

work harness cycle execution depends on a few invariants: work happens on the required branch, state moves through canonical lanes, observer runs before commit, and each completed cycle produces exactly one commit.

## Intended Audience

- operators
- any stage agent that can modify repository state

## Preconditions

- a cycle is being executed or closed

## Main Flow

Treat these as invariant rules:

- stage launch requires the expected branch, usually `dev`
- backlog state follows canonical lanes and does not bypass `qa`
- no stage-level commit is made before observer closure
- each completed cycle has one commit named `cycle-<cycle-id>`
- `dev -> main` promotion is PR-only with human approval

If any invariant would be broken, stop and repair the process rather than continuing with partial compliance.

## Details

These rules keep work harness observable and auditable. They ensure that a cycle artifact set, observer report, and final commit describe the same unit of work.

## Failure Modes

- feature work on the wrong branch
- direct `active -> done` transitions
- multiple commits for one cycle
- observer report missing from the final cycle closure

## References

- [[CYCLE_NAVIGATION]]
- [[STAGE_CONTRACTS]]
- `products/work-harness/DEVELOPMENT_CYCLE.md`
- `products/work-harness/HUMANS.md`
