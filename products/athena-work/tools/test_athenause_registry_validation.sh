#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
validator="$root_dir/tools/platform/validate_athenause_registry.sh"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

valid_registry="$tmp_dir/valid-registry.yaml"
invalid_registry="$tmp_dir/invalid-registry.yaml"

cat > "$valid_registry" <<'EOF'
version: 1
tools:
  - id: athena.test.valid
    name: Valid Tool
    description: Valid registry fixture
    tags: [test]
    stage_affinity: [engineering]
    credential: ""
    call:
      type: exec
      command: ./echo-valid
    schema: []
EOF

cat > "$invalid_registry" <<'EOF'
version: 1
tools:
  - id: athena.test.invalid
    name: Invalid Tool
    description: Missing call type fixture
    tags: [test]
    stage_affinity: [engineering]
    credential: ""
    call:
      command: ./echo-invalid
    schema: []
EOF

ATHENA_USE_REGISTRY="$valid_registry" "$validator" >/dev/null
echo "PASS: validator accepts a valid approved registry"

set +e
invalid_output="$(ATHENA_USE_REGISTRY="$invalid_registry" "$validator" 2>&1)"
code=$?
set -e
if [[ $code -ne 0 ]] && grep -Fq "missing call.type" <<<"$invalid_output"; then
  echo "PASS: validator fails fast with actionable error for invalid registry"
else
  echo "FAIL: validator did not reject invalid registry as expected"
  echo "$invalid_output"
  exit 1
fi
