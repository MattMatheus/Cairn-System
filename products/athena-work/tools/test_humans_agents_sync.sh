#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
cd "$root_dir"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "FAIL: repository context required"
  exit 1
fi

# Files that define stage behavior and process rules.
watched=(
  "products/athena-work/DEVELOPMENT_CYCLE.md"
  "products/athena-work/tools/launch_stage.sh"
  "products/athena-work/tools/run_observer_cycle.sh"
  "products/athena-work/operating-system/README.md"
  "products/athena-work/operating-system/observer/README.md"
  "docs/operator/athena-work/process/STAGE_EXIT_GATES.md"
  "docs/operator/athena-work/process/PROGRAM_OPERATING_SYSTEM.md"
  "docs/operator/athena-work/process/README.md"
  "products/athena-work/delivery-backlog/engineering/intake/STORY_TEMPLATE.md"
  "products/athena-work/delivery-backlog/engineering/intake/BUG_TEMPLATE.md"
  "products/athena-work/delivery-backlog/architecture/intake/ARCH_STORY_TEMPLATE.md"
  "products/athena-work/stage-prompts/active/planning-seed-prompt.md"
  "products/athena-work/stage-prompts/active/next-agent-seed-prompt.md"
  "products/athena-work/stage-prompts/active/architect-agent-seed-prompt.md"
  "products/athena-work/stage-prompts/active/qa-agent-seed-prompt.md"
  "products/athena-work/stage-prompts/active/pm-refinement-seed-prompt.md"
  "products/athena-work/stage-prompts/active/cycle-seed-prompt.md"
)

# Files that must be touched when behavior changes.
sync_targets=(
  "products/athena-work/HUMANS.md"
  "products/athena-work/AGENTS.md"
)

changed_files="$(git status --porcelain | sed -E 's/^.. //')"
if [[ -z "$changed_files" ]]; then
  echo "PASS: no working-tree changes"
  echo "Result: PASS"
  exit 0
fi

watched_changed=0
for f in "${watched[@]}"; do
  if printf '%s\n' "$changed_files" | grep -Fxq "$f"; then
    watched_changed=1
    break
  fi
done

if [[ "$watched_changed" -eq 0 ]]; then
  echo "PASS: no workflow-rule files changed"
  echo "Result: PASS"
  exit 0
fi

sync_touched=0
for f in "${sync_targets[@]}"; do
  if printf '%s\n' "$changed_files" | grep -Fxq "$f"; then
    sync_touched=1
    break
  fi
done

if [[ "$sync_touched" -eq 1 ]]; then
  echo "PASS: workflow-rule change includes HUMANS.md/AGENTS.md sync"
  echo "Result: PASS"
  exit 0
fi

echo "FAIL: workflow-rule files changed without updating HUMANS.md or AGENTS.md"
echo "Changed files:"
printf '%s\n' "$changed_files"
echo "Required sync targets: HUMANS.md or AGENTS.md"
echo "Result: FAIL"
exit 1
