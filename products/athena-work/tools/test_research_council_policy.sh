#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

council_doc="$root_dir/knowledge-base/process/RESEARCH_COUNCIL_BASELINE.md"
planning_prompt="$root_dir/stage-prompts/active/planning-seed-prompt.md"
pm_prompt="$root_dir/stage-prompts/active/pm-refinement-seed-prompt.md"
architect_prompt="$root_dir/stage-prompts/active/architect-agent-seed-prompt.md"

doc_test_init

doc_assert_exists "$council_doc" "Research Council baseline doc exists"
doc_assert_contains "$council_doc" "User-Directed, Agent-Executed Updates" "Council doc defines user-directed agent-executed governance"
doc_assert_contains "$council_doc" "ATHENA_DIRECTION_CONFIRMATION_ID" "Council doc requires explicit direction confirmation evidence"

doc_assert_contains "$planning_prompt" "RESEARCH_COUNCIL_BASELINE.md" "Planning prompt references Research Council baseline"
doc_assert_contains "$pm_prompt" "COUNCIL-" "PM prompt consumes council artifacts"
doc_assert_contains "$architect_prompt" "COUNCIL-" "Architect prompt consumes council artifacts"

doc_test_finish
