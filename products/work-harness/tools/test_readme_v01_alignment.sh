#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$script_dir/lib/doc_test_harness.sh"

readme="$root_dir/README.md"
vision="$root_dir/docs/product/work-harness/product/VISION.md"

doc_test_init

doc_assert_exists "$readme" "README exists"
doc_assert_exists "$vision" "Vision document exists"
doc_assert_contains "$readme" "What v0.1 Delivers Today" "README includes current v0.1 scope section"
doc_assert_contains "$readme" "go run ./cmd/state-harness tooling write" "README includes write example"
doc_assert_contains "$readme" "go run ./cmd/state-harness tooling retrieve" "README includes retrieve example"
doc_assert_contains "$readme" "product-research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md" "README links phased plan"
doc_assert_contains "$readme" "docs/product/work-harness/product/VISION.md" "README links preserved vision doc"
doc_assert_contains "$vision" "work harness Product Vision (Long-Term)" "Vision document title present"

for forbidden in "FAISS" "Podman" "SQLite" "embeddings" "cloud"; do
  if grep -qi "$forbidden" "$readme"; then
    echo "FAIL: README should avoid unimplemented term '$forbidden'"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  else
    echo "PASS: README avoids unimplemented term '$forbidden'"
  fi
done

doc_test_finish
