#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/../.." && pwd))"
athena_home_default="$root_dir/.athena"
athena_home="${ATHENA_HOME:-$athena_home_default}"
download_modules=false

usage() {
  cat <<EOF
usage: tools/dev/bootstrap_platform.sh [--download-modules]

Prepares the repo-local Athena runtime area and verifies the local Go toolchain.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --download-modules)
      download_modules=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "error: unknown arg '$1'" >&2
      usage
      exit 1
      ;;
  esac
done

mkdir -p \
  "$athena_home/workspace" \
  "$athena_home/memory/default" \
  "$athena_home/artifacts/releases" \
  "$athena_home/cache" \
  "$athena_home/runs" \
  "$athena_home/config" \
  "$athena_home/bin"

"$root_dir/products/athena-work/tools/check_go_toolchain.sh"

if "$download_modules"; then
  (cd "$root_dir/products/athena-mind" && go mod download)
  (cd "$root_dir/products/athena-use" && go mod download)
fi

cat <<EOF
status: ready
athena_home: $athena_home
memory_root: $athena_home/memory/default
next_steps:
  1. export ATHENA_HOME="$athena_home"
  2. export ATHENA_MEMORY_ROOT="$athena_home/memory/default"
  3. review PLATFORM_QUICKSTART.md
  4. run ./tools/dev/build_release_artifacts.sh --version <label>
EOF
