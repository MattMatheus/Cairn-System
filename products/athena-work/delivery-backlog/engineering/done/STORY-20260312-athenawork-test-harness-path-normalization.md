# Story: Normalize remaining AthenaWork test-harness repo paths

## Metadata
- `id`: STORY-20260312-athenawork-test-harness-path-normalization
- `owner_persona`: Product Manager - AthenaWork
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0004]
- `success_metric`: AthenaWork validation scripts use the unified repo layout consistently instead of mixing historical root-level paths
- `release_checkpoint`: deferred

## Problem Statement
- Canonical doc validation is fixed, but several AthenaWork test scripts still reference pre-merge root-level paths and will drift further if left unnormalized.

## Scope
- In:
  - audit remaining `products/athena-work/tools/*.sh` scripts for historical repo-era path assumptions
  - normalize pathing to current `products/athena-work/` and `docs/` surfaces
  - update any directly affected tests/docs
- Out:
  - broader workflow redesign
  - backlog/state-model changes
  - non-AthenaWork product scripts

## Acceptance Criteria
1. Remaining AthenaWork validation scripts use consistent current-layout path resolution.
2. Canonical validation entrypoints continue to pass after normalization.
3. Historical-path cleanup does not alter stage behavior beyond path correctness.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- Current AthenaWork tool/test harness under `products/athena-work/tools/`
- Unified repo layout under `products/`, `tools/`, `docs/`, and `.athena/`

## Notes
- Keep this behind product stories unless another validation gate is blocked by path drift again.
