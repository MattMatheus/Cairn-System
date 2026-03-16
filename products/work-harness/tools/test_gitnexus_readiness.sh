#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/../../.." && pwd))"
source "$script_dir/lib/doc_test_harness.sh"

howto_doc="$root_dir/docs/operator/work-harness/how-to/GITNEXUS_READINESS.md"
workflow_doc="$root_dir/docs/operator/work-harness/process/OPERATOR_DAILY_WORKFLOW.md"
readiness_script="$script_dir/check_gitnexus_readiness.sh"

doc_test_init

doc_assert_exists "$readiness_script" "GitNexus readiness script exists"
doc_assert_exists "$howto_doc" "GitNexus readiness doc exists"
doc_assert_contains "$howto_doc" 'check_gitnexus_readiness.sh' "How-to references readiness script"
doc_assert_contains "$howto_doc" 'codegraph-cli status' "How-to references codegraph status command"
doc_assert_contains "$workflow_doc" 'check_gitnexus_readiness.sh' "Daily workflow includes GitNexus readiness preflight"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

mkdir -p "$tmp_dir/bin" "$tmp_dir/repos/untrusted/GitNexus/gitnexus/dist/cli"
cat > "$tmp_dir/bin/node" <<'EOF'
#!/usr/bin/env bash
exit 0
EOF
chmod +x "$tmp_dir/bin/node"
touch "$tmp_dir/repos/untrusted/GitNexus/gitnexus/dist/cli/index.js"

output="$(PATH="$tmp_dir/bin:$PATH" "$readiness_script" --root "$tmp_dir" 2>&1 || true)"
if grep -Fq "PASS: GitNexus ready via built local checkout" <<<"$output"; then
  echo "PASS: readiness script accepts built local checkout with node on PATH"
else
  echo "FAIL: readiness script accepts built local checkout with node on PATH"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
fi

doc_test_finish
