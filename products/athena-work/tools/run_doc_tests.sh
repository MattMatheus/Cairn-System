#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

# State-harness baseline validation suite.
"$root_dir/tools/test_no_personal_paths.sh"
"$root_dir/tools/test_launch_stage_workspace_api_adapter.sh"
"$root_dir/tools/test_observer_workspace_api_adapter.sh"
"$root_dir/tools/test_workspace_api_state_machine_v1.sh"
"$root_dir/tools/test_v1_v2_zip_artifact_ingest.sh"
"$root_dir/tools/test_research_council_policy.sh"
