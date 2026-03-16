# work harness Validation Surface Status

This file classifies work harness validation scripts that still exist under `products/work-harness/tools/`.

## Current

These are part of the active validation surface and should stay runnable against the unified repo layout:

- `run_doc_tests.sh`
- `validate_intake_items.sh`
- `test_tool_registry_validation.sh`
- `test_docs_navigation_hardening.sh`
- `test_docs_workspace_linkage_v1.sh`
- `test_founder_operator_workflow.sh`
- `test_go_toolchain_readiness.sh`
- `test_humans_agents_sync.sh`
- `test_human_planning_workbench_v1.sh`
- `test_intake_validation.sh`
- `test_launch_stage_readme_queue.sh`
- `test_launch_stage_workspace_api_adapter.sh`
- `test_local_control_plane_bootstrap.sh`
- `test_markdown_drift_guard_regression_v1.sh`
- `test_no_personal_paths.sh`
- `test_observer_cycle_policy.sh`
- `test_observer_workspace_api_adapter.sh`
- `test_qa_regression_rubric.sh`
- `test_research_council_policy.sh`
- `test_release_launch_authorization_workbench_v1.sh`
- `test_security_nonce_gate.sh`
- `test_stage_exit_pipeline.sh`
- `test_state_transition_checklist.sh`
- `test_v1_v2_zip_artifact_ingest.sh`
- `test_workspace_api_state_machine_v1.sh`
- `test_workspace_entity_relationship_model_v1.sh`
- `test_workspace_ui_read_only_board_v1.sh`

## Historical

These scripts preserve validation intent from the pre-merge repo era, but they still depend on historical artifacts or superseded workflow assumptions and should not be treated as active coverage:

- `test_doc_test_harness_standardization.sh`
- `test_dogfood_scenario_pack_v01.sh`
- `test_kpi_snapshot_baseline_v01.sh`
- `test_kpi_snapshot_delta_post_hardening_v01.sh`
- `test_markdown_sync_worker_and_drift_guard_v1.sh`
- `test_phased_plan_v01_v03.sh`
- `test_program_state_consistency.sh`
- `test_readme_v01_alignment.sh`
- `test_release_checkpoint_bundle_v01.sh`

## Retired

No scripts are physically removed in this checkpoint. If a historical script is rewritten against the active repo surface, move it back into `Current` in this file as part of that story.
