#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/../.." && pwd))"
cairn_home_default="$root_dir/.cairn"
cairn_home="${CAIRN_HOME:-$cairn_home_default}"
download_modules=false

usage() {
  cat <<EOF
usage: tools/dev/bootstrap_platform.sh [--download-modules]

Prepares the repo-local Cairn runtime area and verifies the local Go toolchain.
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
  "$cairn_home/workspace" \
  "$cairn_home/memory/default" \
  "$cairn_home/artifacts/releases" \
  "$cairn_home/cache" \
  "$cairn_home/runs" \
  "$cairn_home/config" \
  "$cairn_home/bin"

"$root_dir/products/work-harness/tools/check_go_toolchain.sh"

if "$download_modules"; then
  (cd "$root_dir/products/memory-cli" && go mod download)
  (cd "$root_dir/products/tool-cli" && go mod download)
fi

cat <<EOF
status: ready
cairn_home: $cairn_home
memory_root: $cairn_home/memory/default
next_steps:
  1. export CAIRN_HOME="$cairn_home"
  2. export CAIRN_MEMORY_ROOT="$cairn_home/memory/default"
  3. review PLATFORM_QUICKSTART.md
  4. run ./tools/dev/build_release_artifacts.sh --version <label>
EOF
