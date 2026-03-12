# AthenaWork Test Classification

## Purpose

This document classifies imported AthenaWork validation scripts into:

- active candidates for AthenaPlatform
- historical references that need rewrite before promotion

## Active Candidates

These scripts align reasonably well with the current platform shape or represent validation intent worth preserving immediately:

- `products/athena-work/tools/test_launch_stage_workspace_api_adapter.sh`
- `products/athena-work/tools/test_observer_workspace_api_adapter.sh`
- `products/athena-work/tools/test_markdown_drift_guard_regression_v1.sh`
- `products/athena-work/tools/test_intake_validation.sh`

These are the best first candidates for adaptation into trusted platform checks.

## Transitional Core Scripts

These are not themselves `test_*` files, but they are part of the active validation path:

- `products/athena-work/tools/run_doc_tests.sh`
- `products/athena-work/tools/check_markdown_drift.sh`
- `products/athena-work/tools/validate_intake_items.sh`

## Historical Or Repo-Shape-Dependent Tests

Most remaining `test_*` scripts should currently be treated as historical reference material.

Reasons include:

- they assume the standalone AthenaWork repository shape
- they depend on files not yet promoted into the active platform surface
- they validate dated handoff or release artifacts rather than the current platform path

Examples:

- `test_release_checkpoint_bundle_v01.sh`
- `test_release_launch_authorization_workbench_v1.sh`
- `test_kpi_snapshot_baseline_v01.sh`
- `test_kpi_snapshot_delta_post_hardening_v01.sh`
- `test_phased_plan_v01_v03.sh`
- `test_workspace_ui_read_only_board_v1.sh`
- `test_workspace_entity_relationship_model_v1.sh`
- `test_workspace_api_state_machine_v1.sh`
- `test_v1_v2_zip_artifact_ingest.sh`
- `test_program_state_consistency.sh`
- `test_stage_exit_pipeline.sh`

## Promotion Strategy

1. Keep a small active set and adapt them to AthenaPlatform paths.
2. Leave historical tests in place but mark them as non-authoritative until rewritten.
3. Prefer a few durable, trusted checks over a large imported test surface with unclear validity.

