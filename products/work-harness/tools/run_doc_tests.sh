#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# State-harness baseline validation suite.
"$script_dir/test_no_personal_paths.sh"
"$script_dir/test_launch_stage_workspace_api_adapter.sh"
"$script_dir/test_observer_workspace_api_adapter.sh"
"$script_dir/test_workspace_api_state_machine_v1.sh"
"$script_dir/test_v1_v2_zip_artifact_ingest.sh"
"$script_dir/test_research_council_policy.sh"
"$script_dir/test_gitnexus_readiness.sh"
