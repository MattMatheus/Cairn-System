# Engineering Handoff: STORY-20260312-work-harness-test-harness-path-normalization

## What Changed
- Normalized current-layout path resolution across the active work harness validation chain, including `run_doc_tests.sh`, `validate_intake_items.sh`, `test_stage_exit_pipeline.sh`, `test_intake_validation.sh`, and several directly related doc-test scripts.
- Fixed `validate_intake_items.sh` to resolve backlog lanes from `products/work-harness/` instead of the repo root, removing a false-positive pass path.
- Updated multiple work harness validation scripts to use product-local harness libraries and current `docs/operator/work-harness/` / `products/work-harness/` paths.
- Logged a follow-on intake item for second-tier historical validation scripts that still assume older artifacts or repo-era layout.

## Why It Changed
- The unified repo migration left several validation scripts with stale path assumptions.
- The most important consequence was that some canonical checks could fail on bad paths or pass without validating the intended files.

## Test Updates Made
- No new standalone test file was added; the change was verified by running the repaired canonical validation entrypoints directly.

## Test Run Results
- `./products/work-harness/tools/validate_intake_items.sh`: pass
- `./products/work-harness/tools/test_intake_validation.sh`: pass
- `./products/work-harness/tools/test_stage_exit_pipeline.sh`: pass
- `./products/work-harness/tools/run_doc_tests.sh`: pass

## Open Risks/Questions
- A second-tier set of older work harness validation scripts still references historical paths or artifacts and is now tracked separately in `STORY-20260312-work-harness-historical-validation-phase2.md`.
- This cycle intentionally prioritized active validation coverage over exhaustive cleanup of every imported historical test script.

## Recommended QA Focus Areas
- Verify `validate_intake_items.sh` now resolves the real intake lanes under `products/work-harness/`.
- Verify the canonical doc-test chain still passes from the repo root.
- Verify the follow-on intake story captures the residual historical-script cleanup scope clearly.

## New Gaps Discovered
- `products/work-harness/delivery-backlog/engineering/intake/STORY-20260312-work-harness-historical-validation-phase2.md`
