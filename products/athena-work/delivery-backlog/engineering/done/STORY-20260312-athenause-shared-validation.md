# Story: Add AthenaUse validation to shared platform checks

## Metadata
- `id`: STORY-20260312-athenause-shared-validation
- `owner_persona`: Product Manager - AthenaWork
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0005, ADR-0006]
- `success_metric`: shared platform checks fail fast when the approved AthenaUse registry becomes invalid
- `release_checkpoint`: required

## Problem Statement
- AthenaUse validation currently works when run directly, but it is not yet part of the shared platform verification path, so registry breakage can bypass normal checks.

## Scope
- In:
  - Identify the shared platform check entrypoint used for repo validation
  - Add AthenaUse registry validation to that path
  - Document or surface failures clearly enough for engineering and QA handoff
- Out:
  - New execution/runtime behavior for tool calls
  - Registry schema redesign
  - Broader CI system redesign

## Acceptance Criteria
1. The shared platform validation path invokes AthenaUse registry validation.
2. A broken approved registry causes the shared validation path to fail with an actionable error.
3. Tests or verification coverage are added or updated for the new validation path.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- Shared platform smoke/check entrypoints
- AthenaUse `validate` command
- Story `STORY-20260312-athenause-registry-expansion`

## Notes
- This stays behind registry expansion in ranking because it is more valuable after the approved tool surface is less trivial.
