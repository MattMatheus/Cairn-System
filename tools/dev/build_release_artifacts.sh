#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/../.." && pwd))"
version=""
out_dir="${ATHENA_RELEASE_OUT_DIR:-$root_dir/.athena/artifacts/releases}"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

usage() {
  cat <<EOF
usage: tools/dev/build_release_artifacts.sh --version <label> [--out-dir <dir>]

Builds AthenaPlatform CLI release artifacts for:
  - darwin/arm64
  - windows/amd64
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      version="${2:-}"
      shift 2
      ;;
    --out-dir)
      out_dir="${2:-}"
      shift 2
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

if [[ -z "$version" ]]; then
  echo "error: --version is required" >&2
  exit 1
fi

artifact_root="$out_dir/$version"
mkdir -p "$artifact_root"

checksum_cmd="shasum -a 256"
if command -v sha256sum >/dev/null 2>&1; then
  checksum_cmd="sha256sum"
fi

build_one() {
  local module_dir="$1"
  local package_path="$2"
  local binary_name="$3"
  local goos="$4"
  local goarch="$5"
  local ext="$6"
  local staged_dir="$tmp_dir/${binary_name}_${goos}_${goarch}"
  local output_name="${binary_name}_${version}_${goos}_${goarch}.zip"
  local built_bin="$staged_dir/${binary_name}${ext}"

  mkdir -p "$staged_dir"
  (
    cd "$module_dir"
    GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -o "$built_bin" "$package_path"
  )
  (
    cd "$staged_dir"
    zip -q "$artifact_root/$output_name" "${binary_name}${ext}"
  )
}

"$root_dir/products/athena-work/tools/check_go_toolchain.sh"

build_one "$root_dir/products/athena-mind" "./cmd/memory-cli" "memory-cli" "darwin" "arm64" ""
build_one "$root_dir/products/athena-mind" "./cmd/memory-cli" "memory-cli" "windows" "amd64" ".exe"
build_one "$root_dir/products/athena-use" "./cmd/use-cli" "use-cli" "darwin" "arm64" ""
build_one "$root_dir/products/athena-use" "./cmd/use-cli" "use-cli" "windows" "amd64" ".exe"

(
  cd "$artifact_root"
  rm -f SHA256SUMS
  for file in ./*.zip; do
    $checksum_cmd "$file"
  done | sed 's# \./#  #; s# \*./#  #' > SHA256SUMS
)

echo "status: built"
echo "artifact_root: $artifact_root"
