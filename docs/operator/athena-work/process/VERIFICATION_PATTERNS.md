---
title: Verification Patterns
description: Decision rules for choosing the minimum correct verification set before AthenaWork work is handed off or closed.
doc_type: procedure
status: draft
domain: process
canonical_sources:
  - go
related:
  - [[HANDOFF_EXPECTATIONS]]
  - [[BRANCH_AND_COMMIT_INVARIANTS]]
---

# Verification Patterns

## Summary

AthenaWork verification is risk-based but never optional. Run the narrowest set of checks that can credibly prove the touched behavior, then include those commands and results in the handoff.

## Intended Audience

- engineering-stage agents
- QA-stage agents
- operators deciding whether a stage can advance

## Preconditions

- code, docs, config, or workflow behavior changed

## Main Flow

Use this order:

1. Run targeted tests for the touched behavior.
2. Run `products/athena-work/tools/run_doc_tests.sh` when docs, prompts, or process behavior changed.
3. Run broader repo validation when the change crosses package or system boundaries.
4. Run `go test ./...` in each touched Go module when Go code or shared Go tooling is affected.
5. Record the exact commands and outcomes in the handoff.

Choose the smallest credible set, but do not skip required gates just because the change feels local.

## Details

Verification should answer:

- did the changed behavior work
- did related behavior regress
- did the work-system docs remain coherent
- is the next stage safe to proceed

If a check cannot be run, the handoff must say why and what risk remains.

## Failure Modes

- claiming "tested" without commands
- running only global tests and missing the touched behavior
- running only local checks when shared interfaces changed
- omitting module-scoped `go test ./...` or doc tests when required by the touched area

## References

- [[HANDOFF_EXPECTATIONS]]
- [[BRANCH_AND_COMMIT_INVARIANTS]]
- `docs/operator/athena-work/process/STAGE_EXIT_GATES.md`
