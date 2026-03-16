# Engineering Handoff: STORY-20260312-athenawork-alpha-release-readiness

## What Changed
- Added the first real release checkpoint bundle at `products/athena-work/operating-system/handoff/RELEASE_BUNDLE_alpha-2026-03-12.md`.
- Updated `generate_launch_authorization_package.sh` to auto-discover the latest dated release bundle, capture its `ship|hold|missing` decision, and block launches when the bundle decision is not `ship`.
- Updated `validate_launch_authorization_package.sh` and `test_release_launch_authorization_workbench_v1.sh` so launch package output must include the bundle decision.
- Added `docs/operator/athena-work/operations/ALPHA_RELEASE_CHECKPOINT.md` and linked it from the active AthenaWork operator guides.
- Corrected the engineering gate/release docs to use module-scoped Go test commands that match the actual repo layout.

## Why It Changed
- AthenaWork was runnable but not credibly releasable: launch authorization depended on a hardcoded historical release bundle filename, and there was no operator-grade alpha checkpoint path for another person to follow.

## Test Updates Made
- Extended launch authorization validation coverage to require `release_bundle_decision`.
- Reused the current launch/workbench and workspace UI tests to verify the live release surface.

## Test Run Results
- `./products/athena-work/tools/run_doc_tests.sh`: pass
- `./products/athena-work/tools/test_workspace_ui_read_only_board_v1.sh`: pass
- `./products/athena-work/tools/test_release_launch_authorization_workbench_v1.sh`: pass
- `cd products/athena-use && go test ./...`: pass
- `cd products/athena-mind && go test ./...`: pass
- `./products/athena-work/tools/generate_launch_authorization_package.sh --root /Users/mattmatheus/Cairn`: pass; only remaining blocker after queue clear is direction confirmation

## Open Risks/Questions
- Human confirmation is still intentionally required before this alpha can be considered shipped.
- The release bundle is adequate for alpha feedback, but it does not imply production deployment readiness.

## Recommended QA Focus Areas
- Verify the generated launch package now points at the dated alpha release bundle and records `release_bundle_decision`.
- Verify the post-story launch blocker set is reduced to explicit human confirmation only.
- Verify the alpha checkpoint doc is sufficient for another operator to reproduce the release-evidence path.
