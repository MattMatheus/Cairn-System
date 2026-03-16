# Alpha Release Checkpoint

Use this path when work harness is ready for external alpha feedback and you need a reproducible release-evidence checkpoint.

## Goal

Produce a release bundle and launch authorization package where the only remaining blocker is explicit human approval.

## Steps

1. Confirm the engineering active queue is empty in `products/work-harness/delivery-backlog/engineering/active/README.md`.
2. Confirm `products/work-harness/delivery-backlog/engineering/intake/` contains only templates.
3. Run the current validation surface:
   - `./products/work-harness/tools/run_doc_tests.sh`
   - `./products/work-harness/tools/test_workspace_ui_read_only_board_v1.sh`
   - `./products/work-harness/tools/test_release_launch_authorization_workbench_v1.sh`
   - `cd products/tool-cli && go test ./...`
   - `cd products/memory-cli && go test ./...`
4. Update or create a dated release bundle at `products/work-harness/operating-system/handoff/RELEASE_BUNDLE_<label>-YYYY-MM-DD.md`.
5. Set the bundle decision explicitly to `ship` or `hold`.
6. Generate launch authorization:
   - `./products/work-harness/tools/generate_launch_authorization_package.sh`
7. Verify only true blockers remain.
   - For a ready-to-approve alpha checkpoint, the expected remaining blocker is human direction confirmation.

## Expected Bundle Contents

- Included stories and bugs
- QA result evidence
- Validation commands and results
- Known risks and rollback direction
- Outcome baseline or expected trend

## Notes

- The launch package automatically selects the latest dated `RELEASE_BUNDLE_*.md` file unless `CAIRN_RELEASE_BUNDLE_PATH` is set.
- A `ship` bundle does not itself ship the release. Human confirmation is still required.
