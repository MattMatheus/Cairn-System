# QA Result: STORY-20260312-platform-bootstrap-and-release-builds

## Verdict
- PASS

## Story
- `products/work-harness/delivery-backlog/engineering/done/STORY-20260312-platform-bootstrap-and-release-builds.md`

## Acceptance Criteria Evidence
- AC1 pass: `tools/dev/bootstrap_platform.sh` now prepares the repo-local `.cairn/` runtime folders and runs `products/work-harness/tools/check_go_toolchain.sh`.
- AC2 pass: `tools/dev/build_release_artifacts.sh` produces deterministic zipped `memory-cli` and `tool-cli` binaries for `darwin/arm64` and `windows/amd64`, plus `SHA256SUMS`.
- AC3 pass: `azure-pipelines.yml` defines an Azure DevOps path that bootstraps the repo, runs work harness doc tests, runs module-scoped Go tests, builds artifacts, and publishes them.

## Test Evidence
- `./tools/dev/bootstrap_platform.sh`: pass
- `./tools/dev/build_release_artifacts.sh --version local-smoke --out-dir /tmp/cairn-release-smoke`: pass
- `./products/work-harness/tools/test_go_toolchain_readiness.sh`: pass
- `./products/work-harness/tools/run_doc_tests.sh`: pass
- `bash -n tools/dev/bootstrap_platform.sh tools/dev/build_release_artifacts.sh products/work-harness/tools/check_go_toolchain.sh`: pass

## Regression Evaluation
- No regression found in the touched setup/build/docs surface.
- The first real Azure pipeline execution remains the main external integration check still outstanding.

## Defects
- None

## Transition Rationale
- Cairn now has a supported bootstrap path and a reproducible binary build contract, which is the right next step after alpha release authorization.
