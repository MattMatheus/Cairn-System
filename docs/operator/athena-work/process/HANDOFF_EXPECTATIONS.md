---
title: Handoff Expectations
description: What AthenaWork handoff packages must contain so the next stage can act without rediscovering implementation context.
doc_type: heuristic
status: draft
domain: process
related:
  - [[STAGE_CONTRACTS]]
  - [[VERIFICATION_PATTERNS]]
---

# Handoff Expectations

## Summary

A handoff package is complete when the next stage can verify, continue, or reject the work without reconstructing the reasoning from git diff alone.

## Intended Audience

- engineering and architect stages preparing handoff
- QA and PM stages consuming handoff artifacts

## Preconditions

- work is moving from one stage to another

## Main Flow

A strong handoff should include:

- what changed
- why it changed
- what was verified
- what remains risky or unresolved
- which files, commands, or artifacts the next stage should inspect first

For engineering handoff, include explicit risks/questions and any new backlog gaps found during implementation. For QA handoff, include enough evidence that acceptance criteria and regression checks can be evaluated directly.

## Details

The handoff is the bridge between stage contracts. It should reduce the next stage's search cost, not just announce completion.

Good handoff signals:

- direct references to changed files
- exact verification commands run
- explicit notes on untested edges or deferred concerns
- traceability back to story or ADR context

## Failure Modes

- "implemented" with no evidence
- missing risk notes
- no indication of what QA should inspect first
- hidden gaps discovered during implementation but not recorded

## References

- [[STAGE_CONTRACTS]]
- [[VERIFICATION_PATTERNS]]
- `docs/operator/athena-work/process/STAGE_EXIT_GATES.md`
