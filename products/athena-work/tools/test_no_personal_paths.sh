#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

target_paths=(
  "AGENTS.md"
  "HUMANS.md"
  "DEVELOPMENT_CYCLE.md"
  "knowledge-base"
  "stage-prompts"
  "operating-system"
)

allowed_globs=(
  "*/delivery-backlog/*/done/*"
)

declare -a scan_files=()

collect_files() {
  local path="$1"
  if [[ -f "$path" ]]; then
    scan_files+=("$path")
    return
  fi
  if [[ ! -d "$path" ]]; then
    return
  fi

  while IFS= read -r file; do
    scan_files+=("$file")
  done < <(find "$path" -type f)
}

should_skip() {
  local file="$1"
  for pattern in "${allowed_globs[@]}"; do
    if [[ "$file" == $pattern ]]; then
      return 0
    fi
  done
  return 1
}

for path in "${target_paths[@]}"; do
  collect_files "$root_dir/$path"
done

failures=0
for file in "${scan_files[@]}"; do
  rel="${file#"$root_dir"/}"
  if should_skip "$rel"; then
    continue
  fi

  if [[ ! "$file" =~ \.(md|sh|yml|yaml|toml|json|go)$ ]]; then
    continue
  fi

  if rg -q '/Users/|C:\\Users\\|/home/[A-Za-z0-9_.-]+/' "$file"; then
    echo "FAIL: personal path detected in $rel"
    failures=$((failures + 1))
  fi
done

if [[ "$failures" -gt 0 ]]; then
  echo "Result: FAIL ($failures files contain personal paths)"
  exit 1
fi

echo "PASS: no personal absolute paths in guarded docs/prompts"
echo "Result: PASS"
