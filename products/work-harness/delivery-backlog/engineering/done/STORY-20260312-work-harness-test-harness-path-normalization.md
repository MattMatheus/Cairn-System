# Story: Normalize remaining work harness test-harness repo paths

## Metadata
- `id`: STORY-20260312-work-harness-test-harness-path-normalization
- `owner_persona`: Product Manager - work harness
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0004]
- `success_metric`: work harness validation scripts use the unified repo layout consistently instead of mixing historical root-level paths
- `release_checkpoint`: deferred

## Problem Statement
- Canonical doc validation is fixed, but several work harness test scripts still reference pre-merge root-level paths and will drift further if left unnormalized.

## Scope
- In:
  - audit remaining `products/work-harness/tools/*.sh` scripts for historical repo-era path assumptions
  - normalize pathing to current `products/work-harness/` and `docs/` surfaces
  - update any directly affected tests/docs
- Out:
  - broader workflow redesign
  - backlog/state-model changes
  - non-work harness product scripts

## Acceptance Criteria
1. Remaining work harness validation scripts use consistent current-layout path resolution.
2. Canonical validation entrypoints continue to pass after normalization.
3. Historical-path cleanup does not alter stage behavior beyond path correctness.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- Current work harness tool/test harness under `products/work-harness/tools/`
- Unified repo layout under `products/`, `tools/`, `docs/`, and `.cairn/`

## Notes
- Keep this behind product stories unless another validation gate is blocked by path drift again.
