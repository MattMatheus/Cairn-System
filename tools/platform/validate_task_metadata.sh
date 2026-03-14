#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TARGET_DIR="$ROOT_DIR/workspace"
MODE="changed"

if [[ "${1:-}" == "--all" ]]; then
  MODE="all"
  shift
elif [[ "${1:-}" == "--changed" ]]; then
  MODE="changed"
  shift
fi

if [[ -n "${1:-}" ]]; then
  TARGET_DIR="$1"
fi

ALLOWED_TASK_TYPES="implementation investigation documentation quality operations decision maintenance"
ALLOWED_OUTPUT_TYPES="code_change design_note qa_evidence runbook_update decision_record docs_update report"

errors=0
checked=0

is_allowed() {
  local value="$1"
  shift
  for option in "$@"; do
    if [[ "$value" == "$option" ]]; then
      return 0
    fi
  done
  return 1
}

collect_files() {
  if [[ "$MODE" == "all" ]]; then
    find "$TARGET_DIR" -type f -name '*.md' -not -path '*/.git/*' -print0
    return
  fi

  if ! command -v git >/dev/null 2>&1 || ! git -C "$ROOT_DIR" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "WARN: git unavailable or repo not detected; falling back to --all mode." >&2
    find "$TARGET_DIR" -type f -name '*.md' -not -path '*/.git/*' -print0
    return
  fi

  {
    git -C "$ROOT_DIR" diff --name-only --diff-filter=ACMRTUXB -- '*.md'
    git -C "$ROOT_DIR" ls-files --others --exclude-standard -- '*.md'
  } | sed '/^$/d' | awk '!seen[$0]++' | while IFS= read -r rel; do
    printf '%s\0' "$ROOT_DIR/$rel"
  done
}

while IFS= read -r -d '' file; do
  if ! head -n 1 "$file" | grep -qx -- '---'; then
    continue
  fi

  frontmatter="$(awk '
    NR == 1 && $0 == "---" { in_fm = 1; next }
    in_fm && $0 == "---" { exit }
    in_fm { print }
  ' "$file")"

  if [[ -z "$frontmatter" ]]; then
    continue
  fi

  type_value="$(printf '%s\n' "$frontmatter" | awk -F': *' '$1=="type"{print $2; exit}' | tr -d '"')"
  if [[ "$type_value" != "task" ]]; then
    continue
  fi

  checked=$((checked + 1))
  file_errors=0

  for key in task_type output_type acceptance_evidence; do
    if ! printf '%s\n' "$frontmatter" | grep -Eq "^${key}:"; then
      echo "ERROR: $file missing required field '$key'"
      errors=$((errors + 1))
      file_errors=$((file_errors + 1))
    fi
  done

  if [[ "$file_errors" -gt 0 ]]; then
    continue
  fi

  task_type="$(printf '%s\n' "$frontmatter" | awk -F': *' '$1=="task_type"{print $2; exit}' | tr -d '"' | xargs)"
  output_type="$(printf '%s\n' "$frontmatter" | awk -F': *' '$1=="output_type"{print $2; exit}' | tr -d '"' | xargs)"

  if ! is_allowed "$task_type" $ALLOWED_TASK_TYPES; then
    echo "ERROR: $file has invalid task_type '$task_type'"
    echo "       allowed: $ALLOWED_TASK_TYPES"
    errors=$((errors + 1))
  fi

  if ! is_allowed "$output_type" $ALLOWED_OUTPUT_TYPES; then
    echo "ERROR: $file has invalid output_type '$output_type'"
    echo "       allowed: $ALLOWED_OUTPUT_TYPES"
    errors=$((errors + 1))
  fi
done < <(collect_files)

if [[ "$checked" -eq 0 ]]; then
  if [[ "$MODE" == "changed" ]]; then
    echo "No changed task notes found."
  else
    echo "No task notes found under $TARGET_DIR"
  fi
  exit 0
fi

if [[ "$errors" -gt 0 ]]; then
  echo "Task metadata validation failed: $errors error(s) across $checked task note(s)."
  exit 1
fi

echo "Task metadata validation passed for $checked task note(s)."

