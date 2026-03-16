# Story: Normalize second-tier work harness historical validation scripts

## Metadata
- `id`: STORY-20260312-work-harness-historical-validation-phase2
- `owner_persona`: Product Manager - work harness
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0004]
- `success_metric`: older work harness validation scripts outside the canonical doc-test chain stop assuming pre-merge repo-era paths or missing historical artifacts
- `release_checkpoint`: deferred

## Problem Statement
- The canonical work harness validation entrypoints now use current repo-local paths, but a second-tier set of older validation scripts still references historical locations or artifacts that are no longer part of the active platform surface.

## Scope
- In:
  - audit remaining `products/work-harness/tools/test_*.sh` scripts outside the canonical validation chain
  - normalize current-path assumptions where a current equivalent exists
  - retire, rewrite, or explicitly quarantine checks that depend on missing historical artifacts
- Out:
  - changes to active tool-cli product scope
  - non-work harness product tooling
  - broader workflow redesign

## Acceptance Criteria
1. Remaining non-canonical work harness validation scripts are classified as current, historical, or retired.
2. Scripts kept as current use the unified repo layout consistently.
3. Scripts that depend on missing historical artifacts are no longer silently treated as active validation coverage.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- Current work harness tool/test harness under `products/work-harness/tools/`
- Migration notes in `docs/migration/work-harness-active-surface.md`

## Notes
- This is follow-on cleanup after the canonical validation chain and active intake validator were repaired.
- This cycle established an explicit validation-surface classification, repaired current workbench-facing tests and scripts, and left only clearly historical checks outside the active validation set.
