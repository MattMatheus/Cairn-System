# Release Bundle: alpha-2026-03-12

## Release Identifier Standard
- Use date-based identifiers: `<label>-YYYY-MM-DD`
- Example: `v0.2-beta-2026-02-27`
- File path format: `operating-system/handoff/RELEASE_BUNDLE_<release-id>.md`

## Decision
- `ship`

Rationale: current work harness alpha scope is coherent enough for feedback collection, and the remaining launch blocker should be explicit human approval rather than missing evidence.

## Scope
- Included stories:
  - `STORY-20260312-toolcli-registry-expansion`
  - `STORY-20260312-toolcli-shared-validation`
  - `STORY-20260312-toolcli-context-schema-output`
  - `STORY-20260312-toolcli-registry-parameter-metadata`
  - `STORY-20260312-work-harness-test-harness-path-normalization`
  - `STORY-20260312-work-harness-historical-validation-phase2`
- Included bugs:
  - `BUG-20260312-run-doc-tests-root-path`
- Excluded deferred items:
  - broader historical test cleanup outside the active release path
  - production deployment automation beyond current local/operator scope

## Evidence
- QA result artifacts:
  - `products/work-harness/delivery-backlog/engineering/done/QA-RESULT-STORY-20260312-toolcli-registry-expansion.md`
  - `products/work-harness/delivery-backlog/engineering/done/QA-RESULT-STORY-20260312-toolcli-shared-validation.md`
  - `products/work-harness/delivery-backlog/engineering/done/QA-RESULT-STORY-20260312-toolcli-context-schema-output.md`
  - `products/work-harness/delivery-backlog/engineering/done/QA-RESULT-STORY-20260312-toolcli-registry-parameter-metadata.md`
  - `products/work-harness/delivery-backlog/engineering/done/QA-RESULT-STORY-20260312-work-harness-test-harness-path-normalization.md`
  - `products/work-harness/delivery-backlog/engineering/done/QA-RESULT-STORY-20260312-work-harness-historical-validation-phase2.md`
- Validation commands/results:
  - `./products/work-harness/tools/run_doc_tests.sh`: pass
  - `./products/work-harness/tools/test_workspace_ui_read_only_board_v1.sh`: pass
  - `./products/work-harness/tools/test_release_launch_authorization_workbench_v1.sh`: pass
  - `cd products/tool-cli && go test ./...`: pass
  - `cd products/memory-cli && go test ./...`: pass
- Program board snapshot reference:
  - `products/work-harness/delivery-backlog/engineering/active/README.md` shows no active execution queue before release authorization

## Risk and Rollback
- Known risks:
  - release validation is strong for local/operator flows but not yet backed by production deployment automation
  - several historical work harness tests remain intentionally outside the active release path
  - external alpha users may encounter onboarding friction outside the currently documented operator path
- Rollback direction:
  - treat this alpha as a feedback checkpoint only
  - if issues surface, revert to the previous `main` checkpoint and mark the next release bundle `hold` until the gap is corrected

## Outcome Signals
- Baseline metric snapshot:
  - launch authorization now reflects true blockers rather than stale path drift
  - tool-cli context and registry metadata are present in stage launch output
- Expected trend direction:
  - faster external feedback on operator usability and release evidence gaps
  - fewer false-negative release blockers caused by repo-shape drift

## Notes
- Human approval is still required via explicit direction confirmation before this checkpoint is actually shipped.
