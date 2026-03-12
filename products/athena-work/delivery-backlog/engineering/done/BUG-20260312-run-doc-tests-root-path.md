# Bug: run_doc_tests.sh references non-existent root-level test scripts

## Metadata
- `id`: BUG-20260312-run-doc-tests-root-path
- `priority`: P1
- `reported_by`: QA Engineer - Iris
- `source_story`: products/athena-work/delivery-backlog/architecture/active/ARCH-20260312-athenause-tool-interface-spec.md
- `status`: done
- `phase`: v0.1
- `adr_refs`: [ADR-0007]
- `impact_metric`: canonical doc-test entrypoint currently exits before running any doc validation, blocking architecture validation cycles

## Priority Definitions
- `P0`: release-blocking, data loss/corruption, security-critical
- `P1`: major functional regression, acceptance criteria blocked
- `P2`: moderate defect with workaround
- `P3`: minor defect, polish, or low-impact inconsistency

## Summary
`products/athena-work/tools/run_doc_tests.sh` currently calls test scripts through `$root_dir/tools/...`, but those scripts live under `products/athena-work/tools/`. As a result, the canonical documentation validation command fails immediately before executing any checks.

## Expected Behavior
- `./products/athena-work/tools/run_doc_tests.sh` runs the intended AthenaWork doc/state-harness test suite from the repo root.

## Actual Behavior
- The script exits immediately with `No such file or directory` for `/Users/mattmatheus/AthenaPlatform/tools/test_no_personal_paths.sh`.

## Reproduction Steps
1. From the repo root, run `./products/athena-work/tools/run_doc_tests.sh`.
2. Observe the script attempt to execute `$root_dir/tools/test_no_personal_paths.sh`.
3. See the command fail because the referenced file does not exist at that path.

## Evidence
- `./products/athena-work/tools/run_doc_tests.sh: line 8: /Users/mattmatheus/AthenaPlatform/tools/test_no_personal_paths.sh: No such file or directory`

## Suggested Fix Direction (Optional)
- Fixed in the same bootstrap sequence by repointing the canonical doc-test entrypoint and its dependent AthenaWork tests to the current product-local paths, then rerunning `./products/athena-work/tools/run_doc_tests.sh` successfully.
