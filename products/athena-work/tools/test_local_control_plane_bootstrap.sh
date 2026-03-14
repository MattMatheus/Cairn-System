#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$script_dir/lib/doc_test_harness.sh"

quickstart="$root_dir/docs/operator/athena-work/operations/LOCAL_CONTROL_PLANE_QUICKSTART.md"
reset_script="$script_dir/workspace_reset.sh"

doc_test_init

doc_assert_exists "$quickstart" "Local control plane quickstart exists"
doc_assert_contains "$quickstart" 'podman compose -f products/athena-work/operations/docker-compose.local.yml up -d --build' "Quickstart documents startup command"
doc_assert_contains "$quickstart" 'podman compose -f products/athena-work/operations/docker-compose.local.yml down' "Quickstart documents teardown command"
doc_assert_contains "$quickstart" '/api/v1/model/respond' "Quickstart documents model response endpoint check"
doc_assert_contains "$quickstart" './products/athena-work/tools/workspace_reset.sh' "Quickstart documents reset command"

doc_assert_exists "$reset_script" "Workspace reset script exists"
doc_assert_contains "$reset_script" 'products/athena-work/operations/docker-compose.local.yml down -v --remove-orphans' "Reset script removes compose volumes"

doc_test_finish
