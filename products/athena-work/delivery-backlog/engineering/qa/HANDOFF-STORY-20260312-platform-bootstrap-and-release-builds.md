# Engineering Handoff: STORY-20260312-platform-bootstrap-and-release-builds

## What Changed
- Added `tools/dev/bootstrap_platform.sh` to create the repo-local `.athena/` runtime shape and verify the shared Go toolchain.
- Added `tools/dev/build_release_artifacts.sh` to build zipped `memory-cli` and `use-cli` artifacts plus `SHA256SUMS` for `darwin/arm64` and `windows/amd64`.
- Added `azure-pipelines.yml` so Azure DevOps runs bootstrap, AthenaWork doc tests, module-scoped Go tests, and publishes release artifacts.
- Fixed `products/athena-work/tools/check_go_toolchain.sh` and its docs/tests to use the real product module layout instead of a nonexistent root `go.mod`.
- Updated quickstart, installation, binary, and verification docs to point at the new bootstrap/build path.

## Why It Changed
- Alpha release readiness improved the evidence path, but setup and binary distribution were still ad hoc.
- Without a supported local bootstrap command and CI build contract, every new user or release would continue to reinvent the same setup steps.

## Test Updates Made
- Updated the Go toolchain readiness doc test to validate module-scoped Go commands.
- Reused AthenaWork doc tests for workflow/doc consistency.

## Test Run Results
- `./tools/dev/bootstrap_platform.sh`: pass
- `./tools/dev/build_release_artifacts.sh --version local-smoke --out-dir /tmp/athena-release-smoke`: pass
- `./products/athena-work/tools/test_go_toolchain_readiness.sh`: pass
- `./products/athena-work/tools/run_doc_tests.sh`: pass
- `bash -n tools/dev/bootstrap_platform.sh tools/dev/build_release_artifacts.sh products/athena-work/tools/check_go_toolchain.sh`: pass

## Open Risks/Questions
- The Azure pipeline config is added but not executed in this local session, so the first live Azure run is still the real integration proof.
- The current binary scope is intentionally narrow and does not yet cover Linux or installer packaging.

## Recommended QA Focus Areas
- Verify the generated artifact names and checksum file match the documented distribution contract.
- Verify the Azure pipeline uses the intended release artifact staging directory and artifact name.
- Verify bootstrap output is clear enough for a first-time operator to follow without extra context.
