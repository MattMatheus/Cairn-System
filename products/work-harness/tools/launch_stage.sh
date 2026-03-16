#!/usr/bin/env bash
set -euo pipefail

stage="${1:-engineering}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
platform_root="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/../../.." && pwd))"
product_root="$(cd "$script_dir/.." && pwd)"
active_dir="$product_root/delivery-backlog/engineering/active"
active_readme="$active_dir/README.md"
arch_active_dir="$product_root/delivery-backlog/architecture/active"
arch_active_readme="$arch_active_dir/README.md"
required_branch="${CAIRN_REQUIRED_BRANCH:-dev}"
source "$product_root/tools/lib/workspace_api_adapter.sh"

emit_tool_context() {
  local stage_name="$1"
  local use_product_root="$platform_root/products/tool-cli"
  local output=""
  local err=0

  if [ -n "${USE_CLI_BIN:-}" ]; then
    if ! output="$("$USE_CLI_BIN" context --stage "$stage_name" 2>&1)"; then
      err=$?
    fi
  elif [ -f "$use_product_root/go.mod" ] && [ -f "$use_product_root/cmd/tool-cli/main.go" ]; then
    if ! output="$(cd "$use_product_root" && go run ./cmd/tool-cli context --stage "$stage_name" 2>&1)"; then
      err=$?
    fi
  elif command -v tool-cli >/dev/null 2>&1; then
    if ! output="$(tool-cli context --stage "$stage_name" 2>&1)"; then
      err=$?
    fi
  else
    echo "warning: tool context skipped; tool-cli unavailable" >&2
    return 0
  fi

  if [ "$err" -ne 0 ]; then
    echo "warning: tool context skipped; tool-cli failed: $output" >&2
    return 0
  fi

  printf '%s\n' "$output"
}

if ! git -C "$platform_root" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "abort: not a git repository at $platform_root" >&2
  exit 1
fi

current_branch="$(git -C "$platform_root" branch --show-current)"
if [ "$current_branch" != "$required_branch" ]; then
  echo "abort: active branch is '$current_branch'; expected '$required_branch'" >&2
  exit 1
fi

select_top_story_from_lane() {
  local lane_dir="$1"
  local lane_readme="$2"
  local candidate

  if [ -f "$lane_readme" ]; then
    while IFS= read -r candidate; do
      [ -z "$candidate" ] && continue
      if [ "${candidate#/}" = "$candidate" ]; then
        if [[ "$candidate" == */* ]]; then
          candidate="$platform_root/$candidate"
        else
          candidate="$lane_dir/$candidate"
        fi
      fi
      if [ -f "$candidate" ]; then
        echo "$candidate"
        return 0
      fi
    done < <(sed -En 's/^[[:space:]]*[0-9]+\.[[:space:]]*`([^`]+)`.*/\1/p' "$lane_readme")
  fi

  find "$lane_dir" -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | sort | head -n1
}

case "$stage" in
  engineering)
    workspace_api_handle_direction_confirmation "launch-engineering" "${CAIRN_CYCLE_ID:-launch-engineering}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-engineering}"
    workspace_api_handle_research_comm_exception "launch-engineering" "${CAIRN_CYCLE_ID:-launch-engineering}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-engineering}"
    top_story="$(select_top_story_from_lane "$active_dir" "$active_readme" || true)"
    if [ -z "${top_story:-}" ]; then
      echo "no stories"
      exit 0
    fi

    if [ "${top_story#/}" = "$top_story" ]; then
      if [[ "$top_story" == */* ]]; then
        top_story="$platform_root/$top_story"
      else
        top_story="$active_dir/$top_story"
      fi
    fi

    if [ ! -f "$top_story" ]; then
      echo "abort: top active story not found at '$top_story'" >&2
      exit 1
    fi

    rel_story="${top_story#"$platform_root"/}"
    cat <<EOF
launch: stage-prompts/active/next-agent-seed-prompt.md
cycle: engineering
story: $rel_story
checklist:
  1) read story and clarify open questions
  2) implement required changes
  3) update tests
  4) run tests (must pass)
  5) prepare handoff package
  6) move story to delivery-backlog/engineering/qa
  7) do not commit yet (cycle-level commit after observer step)
EOF
    emit_tool_context "engineering"
    workspace_api_emit_status "launch-engineering"
    ;;
  qa)
    workspace_api_handle_direction_confirmation "launch-qa" "${CAIRN_CYCLE_ID:-launch-qa}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-qa}"
    workspace_api_handle_research_comm_exception "launch-qa" "${CAIRN_CYCLE_ID:-launch-qa}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-qa}"
    cat <<EOF
launch: stage-prompts/active/qa-agent-seed-prompt.md
cycle: qa
checklist:
  1) review story in delivery-backlog/engineering/qa against acceptance criteria
  2) validate tests/regression risk
  3) file defects in delivery-backlog/engineering/intake with P0-P3 if found
  4) move story to delivery-backlog/engineering/done or delivery-backlog/engineering/active
  5) run observer: tools/run_observer_cycle.sh --cycle-id <story-id>
  6) commit once for the full cycle with message: cycle-<cycle-id>
EOF
    emit_tool_context "qa"
    workspace_api_emit_status "launch-qa"
    ;;
  pm)
    workspace_api_handle_direction_confirmation "launch-pm" "${CAIRN_CYCLE_ID:-launch-pm}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-pm}"
    workspace_api_handle_research_comm_exception "launch-pm" "${CAIRN_CYCLE_ID:-launch-pm}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-pm}"
    cat <<EOF
launch: stage-prompts/active/pm-refinement-seed-prompt.md
cycle: pm
checklist:
  1) review/refine items from delivery-backlog/engineering/intake
  2) rank and move selected items to delivery-backlog/engineering/active
  3) update delivery-backlog/engineering/active/README.md sequence
  4) update engineering directive only if needed
  5) run observer: tools/run_observer_cycle.sh --cycle-id PM-<date>-<slug>
  6) commit once for the full cycle with message: cycle-<cycle-id>
EOF
    emit_tool_context "pm"
    workspace_api_emit_status "launch-pm"
    ;;
  planning)
    workspace_api_handle_direction_confirmation "launch-planning" "${CAIRN_CYCLE_ID:-launch-planning}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-planning}"
    workspace_api_handle_research_comm_exception "launch-planning" "${CAIRN_CYCLE_ID:-launch-planning}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-planning}"
    cat <<EOF
launch: stage-prompts/active/planning-seed-prompt.md
cycle: planning
checklist:
  1) run an interactive idea-generation session with the human operator
  2) capture structured notes in workspace/research/planning/sessions using the planning template
  3) convert session output into intake stories (engineering and/or architecture) using canonical templates
  4) recommend next stage: architect (for decisions) and/or pm (for prioritization)
  5) run observer: tools/run_observer_cycle.sh --cycle-id <plan-id>
  6) commit once for the full cycle with message: cycle-<cycle-id>
EOF
    emit_tool_context "planning"
    workspace_api_emit_planning_direction_summary
    workspace_api_emit_status "launch-planning"
    ;;
  architect)
    workspace_api_handle_direction_confirmation "launch-architect" "${CAIRN_CYCLE_ID:-launch-architect}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-architect}"
    workspace_api_handle_research_comm_exception "launch-architect" "${CAIRN_CYCLE_ID:-launch-architect}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-architect}"
    top_arch_story="$(select_top_story_from_lane "$arch_active_dir" "$arch_active_readme" || true)"
    if [ -z "${top_arch_story:-}" ]; then
      echo "no stories"
      exit 0
    fi

    rel_arch_story="${top_arch_story#"$platform_root"/}"
    cat <<EOF
launch: stage-prompts/active/architect-agent-seed-prompt.md
cycle: architect
story: $rel_arch_story
checklist:
  1) read story and clarify architecture decision scope
  2) update ADRs/architecture artifacts
  3) run docs validation tests
  4) prepare handoff package
  5) move story to delivery-backlog/architecture/qa
  6) do not commit yet (cycle-level commit after observer step)
EOF
    emit_tool_context "architect"
    workspace_api_emit_status "launch-architect"
    ;;
  cycle)
    workspace_api_handle_direction_confirmation "launch-cycle" "${CAIRN_CYCLE_ID:-launch-cycle}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-cycle}"
    workspace_api_handle_research_comm_exception "launch-cycle" "${CAIRN_CYCLE_ID:-launch-cycle}" "${CAIRN_STORY_ID:-unscoped}" "${CAIRN_SESSION_ID:-launch-cycle}"
    cat <<EOF
launch: stage-prompts/active/cycle-seed-prompt.md
cycle: engineering+qa loop
loop:
  - run: tools/launch_stage.sh engineering
  - if output is "no stories": stop
  - execute engineering story cycle
  - run: tools/launch_stage.sh qa
  - execute QA cycle
  - run observer: tools/run_observer_cycle.sh --cycle-id <story-id>
  - commit once: cycle-<cycle-id>
  - repeat until active backlog is drained
EOF
    emit_tool_context "cycle"
    workspace_api_emit_status "launch-cycle"
    ;;
  *)
    echo "usage: products/work-harness/tools/launch_stage.sh [engineering|qa|pm|planning|architect|cycle]" >&2
    exit 1
    ;;
esac
