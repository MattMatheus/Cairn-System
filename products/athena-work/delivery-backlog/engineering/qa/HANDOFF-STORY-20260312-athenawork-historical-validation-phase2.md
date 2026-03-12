# Engineering Handoff: STORY-20260312-athenawork-historical-validation-phase2

## What Changed
- Classified the remaining second-tier AthenaWork validation surface in `products/athena-work/tools/VALIDATION_SURFACE_STATUS.md` so current and historical scripts are explicit.
- Repaired current operator-facing scripts and prompts that still referenced historical paths, including planning session guidance, local control plane docs indexes, launch authorization package generation, markdown sync docs paths, and docs-site publishing inputs.
- Restored runnable coverage for current workbench tests by updating `test_workspace_ui_read_only_board_v1.sh`, `test_release_launch_authorization_workbench_v1.sh`, and `test_human_planning_workbench_v1.sh` to the unified repo layout.
- Added a repo-local planning session template under `workspace/research/planning/PLANNING_SESSION_TEMPLATE.md` and redirected planning exports to the workspace research area.

## Why It Changed
- A residual set of AthenaWork scripts still mixed active repo surfaces with pre-merge `product-research/*`, root `tools/lib`, missing compose files, or missing website/docs inputs.
- That drift created false failures in live operator checks, especially the launch authorization package generator's workspace UI gate.

## Test Updates Made
- No new standalone harness was added.
- Existing workbench and launch-validation tests were repaired so they can again serve as current coverage.

## Test Run Results
- `./products/athena-work/tools/build_docs_site.sh /tmp/athena-site-phase2`: pass
- `./products/athena-work/tools/test_workspace_ui_read_only_board_v1.sh`: pass
- `./products/athena-work/tools/test_release_launch_authorization_workbench_v1.sh`: pass
- `./products/athena-work/tools/test_human_planning_workbench_v1.sh`: pass
- `./products/athena-work/tools/test_launch_stage_workspace_api_adapter.sh`: pass
- `./products/athena-work/tools/run_doc_tests.sh`: pass
- `./products/athena-work/tools/validate_launch_authorization_package.sh --package <generated>`: pass

## Open Risks/Questions
- Several intentionally historical tests still preserve older validation intent and remain documented as historical rather than rewritten in this cycle.
- Launch authorization stays `blocked` in the current repo state for real reasons: non-empty intake/QA, missing explicit direction confirmation, and no release bundle file.

## Recommended QA Focus Areas
- Verify the launch authorization package no longer reports a false workspace UI doc-test blocker.
- Verify planning/workbench export now targets `workspace/research/planning/sessions`.
- Verify the classification document matches the scripts that remain outside the canonical validation chain.
